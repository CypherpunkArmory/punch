package tunnel

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cypherpunkarmory/punch/utilities"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	defer remote.Close()
	chDone := make(chan bool)
	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

func privateKeyFile(path string) ssh.AuthMethod {
	path = utilities.FixFilePath(path)
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH key file %s", path))
		return nil
	}
	block, rest := pem.Decode(buffer)
	if len(rest) > 0 {
		return nil
	}
	if !x509.IsEncryptedPEMBlock(block) {
		key, err := ssh.ParsePrivateKey(buffer)
		if err != nil {
			fmt.Println(err)
			log.Fatalln(fmt.Sprintf("Cannot parse SSH key file %s", path))
			return nil
		}
		return ssh.PublicKeys(key)
	}
	fmt.Println("Your password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	key, err := ssh.ParsePrivateKeyWithPassphrase(buffer, bytePassword)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(fmt.Sprintf("Cannot parse SSH key file %s", path))
		return nil
	}
	return ssh.PublicKeys(key)

}

func StartReverseTunnel(tunnelConfig *TunnelConfig) {
	sshPort, _ := strconv.Atoi(tunnelConfig.TunnelEndpoint.SSHPort)
	remoteEndpointPort := 3000
	//TODO: fix when creating punch it
	if tunnelConfig.EndpointType[0] == "https" {
		remoteEndpointPort = 3001
	}
	// local service to be forwarded
	var localEndpoint = Endpoint{
		Host: "0.0.0.0",
		Port: tunnelConfig.LocalPort,
	}
	var jumpServerEndpoint = Endpoint{
		Host: "api.holepunch.io",
		Port: 22,
	}
	// remote SSH server
	var serverEndpoint = Endpoint{
		Host: tunnelConfig.TunnelEndpoint.IPAddress,
		Port: sshPort,
	}

	// remote forwarding port (on remote SSH server network)
	var remoteEndpoint = Endpoint{
		Host: "localhost",
		Port: remoteEndpointPort,
	}
	// refer to https://godoc.org/golang.org/x/crypto/ssh for other authentication types
	sshConfig := &ssh.ClientConfig{
		// SSH connection username
		User: "punch",
		Auth: []ssh.AuthMethod{
			privateKeyFile(tunnelConfig.PrivateKeyPath),
			ssh.Password(""),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         0,
	}

	jumpConn, err := ssh.Dial("tcp", jumpServerEndpoint.String(), sshConfig)
	if err != nil {
		tunnelConfig.RestApi.DeleteTunnelAPI(tunnelConfig.Subdomain)
		log.Fatalln(fmt.Printf("Dial INTO jump server error: %s", err))
		os.Exit(1)
	}
	defer jumpConn.Close()
	// Connect to SSH remote server using serverEndpoint
	var serverConn net.Conn
	serverConn, err = jumpConn.Dial("tcp", serverEndpoint.String())
	if err != nil {
		fmt.Println("Failed to connect. Trying again in 10 seconds")
		time.Sleep(10 * time.Second)
		serverConn, err = jumpConn.Dial("tcp", serverEndpoint.String())
		if err != nil {
			tunnelConfig.RestApi.DeleteTunnelAPI(tunnelConfig.Subdomain)
			log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
			os.Exit(1)
		}
	}
	defer serverConn.Close()
	ncc, chans, reqs, err := ssh.NewClientConn(serverConn, serverEndpoint.String(), sshConfig)
	if err != nil {
		tunnelConfig.RestApi.DeleteTunnelAPI(tunnelConfig.Subdomain)
		log.Fatal(err)
		os.Exit(10)
	}

	sClient := ssh.NewClient(ncc, chans, reqs)
	// Listen on remote server port
	listener, err := sClient.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		tunnelConfig.RestApi.DeleteTunnelAPI(tunnelConfig.Subdomain)
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
		os.Exit(1)
	}
	defer listener.Close()

	// This catches CTRL C and closes the ssh
	c := make(chan os.Signal)
	signal.Notify(c,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP,  // "terminal is disconnected"
	)

	go func() {
		<-c
		tunnelConfig.RestApi.StartSession(tunnelConfig.RestApi.RefreshToken)
		tunnelConfig.RestApi.DeleteTunnelAPI(tunnelConfig.Subdomain)
		listener.Close()
		os.Exit(0)
	}()

	fmt.Printf("Now forwarding localhost:%d to %s://%s.%s\n", tunnelConfig.LocalPort, tunnelConfig.EndpointType, tunnelConfig.Subdomain, tunnelConfig.EndpointUrl)
	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localEndpoint.String())
		if err != nil {
		}

		client, err := listener.Accept()
		if err != nil {
		}
		if err == nil && client != nil && local != nil {
			go handleClient(client, local)
		}
	}

}
