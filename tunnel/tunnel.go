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
	"strconv"
	"sync"
	"syscall"
	"time"

	rollbar "github.com/rollbar/rollbar-go"
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
		_, err := io.Copy(client, remote)
		if err != nil {
			//chDone <- false
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			//chDone <- false
		}
		chDone <- true
	}()
	<-chDone
}

func privateKeyFile(path string) (ssh.AuthMethod, error) {
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
	if wg != nil {
		defer wg.Done()
	}
	listener, err := createTunnel(tunnelConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listener.Close()
	var localEndpoint = endpoint{
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
		err := tunnelConfig.RestAPI.StartSession(tunnelConfig.RestAPI.RefreshToken)
		if err != nil {
			fmt.Println("Could not delete tunnel. Use punch cleanup " + tunnelConfig.Subdomain)
			return
		}
		err = tunnelConfig.RestAPI.DeleteTunnelAPI(tunnelConfig.Subdomain)
		if err != nil {
			fmt.Println("Could not delete tunnel. Use punch cleanup " + tunnelConfig.Subdomain)
		}
		os.Exit(0)
	}()

	fmt.Printf("Now forwarding localhost:%d to %s://%s.%s\n",
		tunnelConfig.LocalPort, tunnelConfig.EndpointType, tunnelConfig.Subdomain, tunnelConfig.EndpointURL)
	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, errLocal := net.Dial("tcp", localEndpoint.String())
		client, errClient := listener.Accept()
		if errLocal == nil && errClient == nil && client != nil && local != nil {
			go handleClient(client, local)
		}
	}

}

func createTunnel(tunnelConfig *Config) (net.Listener, error) {
	var listener net.Listener
	sshPort, _ := strconv.Atoi(tunnelConfig.TunnelEndpoint.SSHPort)
	remoteEndpointPort := 3000
	if tunnelConfig.EndpointType == "https" {
		remoteEndpointPort = 3001
	}
	var jumpServerEndpoint = endpoint{
		Host: "api.holepunch.io",
		Port: 22,
	}
	// remote SSH server
	var serverEndpoint = endpoint{
		Host: tunnelConfig.TunnelEndpoint.IPAddress,
		Port: sshPort,
	}

	// remote forwarding port (on remote SSH server network)
	var remoteEndpoint = endpoint{
		Host: "localhost",
		Port: remoteEndpointPort,
	}

	privateKey, err := privateKeyFile(tunnelConfig.PrivateKeyPath)
	if err != nil {
		errDelete := tunnelConfig.RestAPI.DeleteTunnelAPI(tunnelConfig.Subdomain)
		if errDelete != nil {
			fmt.Println("Could not delete tunnel. Use punch cleanup " + tunnelConfig.Subdomain)
		}
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

	jumpConn, err := ssh.Dial("tcp", jumpServerEndpoint.String(), sshConfig)
	if err != nil {
		errDelete := tunnelConfig.RestAPI.DeleteTunnelAPI(tunnelConfig.Subdomain)
		if errDelete != nil {
			fmt.Println("Could not delete tunnel. Use punch cleanup " + tunnelConfig.Subdomain)
		}
		fmt.Printf("dial INTO jump server error: %s", err)
		return listener, err
	}
	// Connect to SSH remote server using serverEndpoint
	var serverConn net.Conn
	serverConn, err = jumpConn.Dial("tcp", serverEndpoint.String())
	if err != nil {
		rollbar.Message("error", "SSH fail(Jump to Remote): "+err.Error())
		rollbar.Wait()
		fmt.Println("Failed to connect. Trying again in 10 seconds")
		time.Sleep(10 * time.Second)
		serverConn, err = jumpConn.Dial("tcp", serverEndpoint.String())
		if err != nil {
			rollbar.Message("error", "SSH failed twice(Jump to Remote): "+err.Error())
			rollbar.Wait()
			errDelete := tunnelConfig.RestAPI.DeleteTunnelAPI(tunnelConfig.Subdomain)
			if errDelete != nil {
				fmt.Println("Could not delete tunnel. Use punch cleanup " + tunnelConfig.Subdomain)
			}
			fmt.Printf("dial INTO remote server error: %s", err)
			return listener, err
		}
	}
	ncc, chans, reqs, err := ssh.NewClientConn(serverConn, serverEndpoint.String(), sshConfig)
	if err != nil {
		errDelete := tunnelConfig.RestAPI.DeleteTunnelAPI(tunnelConfig.Subdomain)
		if errDelete != nil {
			fmt.Println("Could not delete tunnel. Use punch cleanup " + tunnelConfig.Subdomain)
		}
		return listener, err
	}

	sClient := ssh.NewClient(ncc, chans, reqs)
	// Listen on remote server port
	listener, err = sClient.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		errDelete := tunnelConfig.RestAPI.DeleteTunnelAPI(tunnelConfig.Subdomain)
		if errDelete != nil {
			fmt.Println("Could not delete tunnel. Use punch cleanup " + tunnelConfig.Subdomain)
		}
		fmt.Printf("listen open port ON remote server error: %s", err)
		return listener, err
	}
	return listener, nil
}
