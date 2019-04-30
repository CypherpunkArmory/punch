package tunnel

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cypherpunkarmory/punch/backoff"
	log "github.com/sirupsen/logrus"
	"github.com/tj/go-spin"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client io.ReadWriteCloser, remote io.ReadWriteCloser) {
	defer client.Close()
	defer remote.Close()

	chDone := make(chan bool)
	// Start remote -> local data transfer
	go func() {
		amt, err := io.Copy(client, remote)
		if err != nil {
			log.Debugf("Copy Error: %s ", err)
		}
		log.Debugf("Local -> Remote (%d bytes)", amt)
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		amt, err := io.Copy(remote, client)
		if err != nil {
			log.Debugf("Copy Error %s ", err)
		}
		log.Debugf("Local <- Remote (%d bytes)", amt)
		chDone <- true
	}()
	<-chDone
}

func privateKeyFile(path string) (ssh.AuthMethod, error) {
	log.Debugf("Parsing privatekey %s", path)

	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("cannot read SSH key file " + path)
	}
	if len(buffer) == 0 {
		return nil, errors.New("bad key file empty file")
	}

	block, _ := pem.Decode(buffer)
	if block == nil {
		return nil, errors.New("bad key file")
	}

	// Return early if the SSH file is not password protected
	if !x509.IsEncryptedPEMBlock(block) {
		key, errParse := ssh.ParsePrivateKey(buffer)
		if errParse != nil {
			return nil, errors.New("cannot parse SSH key file " + path)
		}
		return ssh.PublicKeys(key), nil
	}

	fmt.Println("Your password: ")

	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, errors.New("could not read your password " + err.Error())
	}

	key, err := ssh.ParsePrivateKeyWithPassphrase(buffer, bytePassword)
	if err != nil {
		return nil, errors.New("cannot parse SSH key file " + path)
	}

	return ssh.PublicKeys(key), nil
}

//StartReverseTunnel Main tunneling function. Handles connections and forwarding
func StartReverseTunnel(tunnelConfig *Config, wg *sync.WaitGroup) {
	defer cleanup(tunnelConfig)

	if wg != nil {
		defer wg.Done()
	}

	listener, err := createTunnel(tunnelConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return
	}

	defer listener.Close()

	var localEndpoint = Endpoint{
		Host: "0.0.0.0",
		Port: tunnelConfig.LocalPort,
	}

	// This catches CTRL C and closes the ssh
	c := make(chan os.Signal)
	signal.Notify(c,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
	)

	go func() {
		<-c
		cleanup(tunnelConfig)
		os.Exit(0)
	}()

	fmt.Printf("\rNow forwarding localhost:%s to %s://%s.%s\n",
		tunnelConfig.LocalPort, tunnelConfig.EndpointType, tunnelConfig.Subdomain, tunnelConfig.EndpointURL.Host)
	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		log.Debugf("Dial to local %s", localEndpoint.String())
		local, errLocal := net.Dial("tcp", localEndpoint.String())
		client, errClient := listener.Accept()
		if errLocal == nil && errClient == nil && client != nil && local != nil {
			go handleClient(client, local)
		}
	}

}

func createTunnel(tunnelConfig *Config) (net.Listener, error) {
	c := make(chan os.Signal)

	lvl, err := log.ParseLevel(tunnelConfig.LogLevel)
	if err != nil {
		log.Errorf("\nLog level %s is not a valid level.", tunnelConfig.LogLevel)
	}

	log.SetLevel(lvl)
	log.Debugf("Debug Logging activated")

	signal.Notify(c,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
	)

	go func() {
		<-c
		log.Debugf("Closing tunnel")
		cleanup(tunnelConfig)
		os.Exit(0)
	}()

	var listener net.Listener

	sshPort := tunnelConfig.TunnelEndpoint.SSHPort
	// FIXME:  This should be a LUT
	remoteEndpointPort := "3000"

	if tunnelConfig.EndpointType == "https" {
		remoteEndpointPort = "3001"
	}

	var jumpServerEndpoint = Endpoint{
		Host: tunnelConfig.ConnectionEndpoint.Hostname(),
		Port: tunnelConfig.ConnectionEndpoint.Port(),
	}

	// remote SSH server
	var serverEndpoint = Endpoint{
		Host: tunnelConfig.TunnelEndpoint.IPAddress,
		Port: sshPort,
	}

	// remote forwarding port (on remote SSH server network)
	var remoteEndpoint = Endpoint{
		Host: "localhost", // localhost here is the remote SSHD daemon container
		Port: remoteEndpointPort,
	}

	privateKey, err := privateKeyFile(tunnelConfig.PrivateKeyPath)
	if err != nil {
		return listener, err
	}

	sshConfig := &ssh.ClientConfig{
		User: "punch",
		Auth: []ssh.AuthMethod{
			privateKey,
			ssh.Password(""),
		},
		//TODO: Maybe fix this. Will be rotating so dont know if possible
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         0,
	}

	log.Debugf("Dial into Jump Server %s", jumpServerEndpoint.String())
	jumpConn, err := ssh.Dial("tcp", jumpServerEndpoint.String(), sshConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error contacting the Holepunch Server.")
		log.Debugf("%s", err)
		return listener, err
	}

	tunnelStarted := false
	go func() {
		s := spin.New()
		for !tunnelStarted {
			fmt.Printf("\rStarting tunnel %s ", s.Next())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	exponentialBackoff := backoff.NewExponentialBackOff()

	// Connect to SSH remote server using serverEndpoint
	var serverConn net.Conn
	for {
		serverConn, err = jumpConn.Dial("tcp", serverEndpoint.String())
		log.Debugf("Dial into SSHD Container %s", serverEndpoint.String())
		if err == nil {
			tunnelStarted = true
			break
		}
		wait := exponentialBackoff.NextBackOff()
		log.Debugf("Backoff Tick %s", wait.String())
		time.Sleep(wait)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(serverConn, serverEndpoint.String(), sshConfig)
	if err != nil {
		return listener, err
	}
	log.Debugf("SSH Connection Established via Jump %s -> %s", jumpServerEndpoint.String(), serverEndpoint.String())

	sClient := ssh.NewClient(ncc, chans, reqs)
	// Listen on remote server port
	listener, err = sClient.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen open port ON remote server error: %s", err)
		return listener, err
	}
	log.Debugf("Open listen port on %s", remoteEndpoint.String())
	return listener, nil
}

func cleanup(config *Config) {
	errSession := config.RestAPI.StartSession(config.RestAPI.RefreshToken)
	errDelete := config.RestAPI.DeleteTunnelAPI(config.Subdomain)
	if errSession != nil || errDelete != nil {
		fmt.Fprintf(os.Stderr, "Could not delete tunnel. Use punch cleanup %s\n", config.Subdomain)
	}
}
