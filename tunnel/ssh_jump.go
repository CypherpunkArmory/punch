package tunnel

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cypherpunkarmory/punch/backoff"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/tj/go-spin"
	"golang.org/x/crypto/ssh"
)

type tunnelStatus struct {
	state string
}

const (
	tunnelError    string = "error"
	tunnelDone     string = "done"
	tunnelStarting string = "starting"
)

func internalEndpoint(endpointType string) (*Endpoint, error) {
	switch endpointType {
	case "http":
		return &Endpoint{
			Host: "localhost", // localhost here is the remote SSHD daemon container
			Port: "3000",
		}, nil
	case "https":
		return &Endpoint{
			Host: "localhost", // localhost here is the remote SSHD daemon container
			Port: "3001",
		}, nil
	case "tcp":
		return &Endpoint{
			Host: "localhost", // localhost here is the remote SSHD daemon container
			Port: "3002",
		}, nil
	default:
		return nil, errors.New("unknown Endpoint Type")
	}
}

func connectToJumpHost(tunnelConfig *Config, semaphore *Semaphore) (*ssh.Client, error) {
	tunnelCreating := tunnelStatus{tunnelStarting}
	createCloseChannel := make(chan os.Signal)
	signal.Notify(createCloseChannel,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
	)
	defer signal.Stop(createCloseChannel)
	go func() {
		<-createCloseChannel
		log.Debugf("Closing tunnel")
		tunnelCreating.state = tunnelError
		for !semaphore.CanRun() {

		}
		defer semaphore.Done()
		cleanup(tunnelConfig)
		os.Exit(0)
	}()
	lvl, err := log.ParseLevel(tunnelConfig.LogLevel)
	if err != nil {
		log.Errorf("\nLog level %s is not a valid level.", tunnelConfig.LogLevel)
	}

	log.SetLevel(lvl)
	log.Debugf("Debug Logging activated")

	var nilJumpConn *ssh.Client
	var jumpServerEndpoint = Endpoint{
		Host: tunnelConfig.ConnectionEndpoint.Hostname(),
		Port: tunnelConfig.ConnectionEndpoint.Port(),
	}
	hostKeyCallBack := dnsHostKeyCallback
	if tunnelConfig.ConnectionEndpoint.Hostname() != "api.holepunch.io" {
		fmt.Println("Ignoring hostkey")
		hostKeyCallBack = ssh.InsecureIgnoreHostKey()
	}
	sshJumpConfig := &ssh.ClientConfig{
		User: "punch",
		Auth: []ssh.AuthMethod{
			ssh.Password(""),
		},
		HostKeyCallback: hostKeyCallBack,
		Timeout:         0,
	}

	log.Debugf("Dial into Jump Server %s", jumpServerEndpoint.String())
	jumpConn, err := ssh.Dial("tcp", jumpServerEndpoint.String(), sshJumpConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error contacting the Holepunch Server.")
		log.Debugf("%s", err)
		return nilJumpConn, err
	}
	log.Debugf("SSH Connection Established to Jump %s", jumpServerEndpoint.String())
	return jumpConn, nil
}

func createTunnel(jumpConn *ssh.Client, tunnelConfig *Config, semaphore *Semaphore) (*ssh.Client, error) {
	tunnelCreating := tunnelStatus{tunnelStarting}
	createCloseChannel := make(chan os.Signal)
	signal.Notify(createCloseChannel,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
	)
	defer signal.Stop(createCloseChannel)
	go func() {
		<-createCloseChannel
		log.Debugf("Closing tunnel")
		tunnelCreating.state = tunnelError
		for !semaphore.CanRun() {

		}
		defer semaphore.Done()
		cleanup(tunnelConfig)
		os.Exit(0)
	}()
	lvl, err := log.ParseLevel(tunnelConfig.LogLevel)
	if err != nil {
		log.Errorf("\nLog level %s is not a valid level.", tunnelConfig.LogLevel)
	}

	log.SetLevel(lvl)
	log.Debugf("Debug Logging activated")

	sshPort := tunnelConfig.TunnelEndpoint.SSHPort

	// remote SSH server
	var serverEndpoint = Endpoint{
		Host: tunnelConfig.TunnelEndpoint.IPAddress,
		Port: sshPort,
	}

	// remote forwarding port (on remote SSH server network)

	privateKey, err := readPrivateKeyFile(tunnelConfig.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	tunnelStartingSpinner(semaphore, &tunnelCreating)
	exponentialBackoff := backoff.NewExponentialBackOff()

	// Connect to SSH remote server using serverEndpoint
	var serverConn net.Conn
	for {
		serverConn, err = jumpConn.Dial("tcp", serverEndpoint.String())
		log.Debugf("Dial into SSHD Container %s", serverEndpoint.String())
		if err == nil {
			tunnelCreating.state = tunnelDone
			break
		}
		wait := exponentialBackoff.NextBackOff()
		log.Debugf("Backoff Tick %s", wait.String())
		time.Sleep(wait)
	}
	sshTunnelConfig := &ssh.ClientConfig{
		User: "punch",
		Auth: []ssh.AuthMethod{
			privateKey,
		},
		//TODO: Maybe fix this. Will be rotating so dont know if possible
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         0,
	}
	ncc, chans, reqs, err := ssh.NewClientConn(serverConn, serverEndpoint.String(), sshTunnelConfig)
	if err != nil {
		return nil, err
	}
	log.Debugf("SSH Connection Established via Jump %s", serverEndpoint.String())

	sClient := ssh.NewClient(ncc, chans, reqs)
	return sClient, nil
}

func tunnelStartingSpinner(lock *Semaphore, tunnelStatus *tunnelStatus) {
	go func() {
		if !lock.CanRun() {
			return
		}
		defer lock.Done()
		s := spin.New()
		for tunnelStatus.state == tunnelStarting {
			fmt.Printf("\rStarting tunnel %s ", s.Next())
			time.Sleep(100 * time.Millisecond)
		}
		if tunnelStatus.state == tunnelDone {
			fmt.Printf("\rStarting tunnel ")
			d := color.New(color.FgGreen, color.Bold)
			d.Printf("✔\n")
		}
	}()
}
