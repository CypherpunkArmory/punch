package tunnel

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
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

//StartReverseTunnel Main tunneling function. Handles connections and forwarding
func StartReverseTunnel(tunnelConfig *Config, wg *sync.WaitGroup, semaphore *Semaphore) {
	defer cleanup(tunnelConfig)

	if wg != nil {
		defer wg.Done()
	}
	listener, err := createTunnel(tunnelConfig, semaphore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	defer listener.Close()

	var localEndpoint = Endpoint{
		Host: "0.0.0.0",
		Port: tunnelConfig.LocalPort,
	}

	// This catches CTRL C and closes the ssh
	startCloseChannel := make(chan os.Signal)
	signal.Notify(startCloseChannel,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
	)

	go func() {
		<-startCloseChannel
		if semaphore.CanRun() {
			cleanup(tunnelConfig)
			os.Exit(0)
		}
	}()

	fmt.Printf("Access your website at %s://%s.%s\n",
		tunnelConfig.EndpointType, tunnelConfig.Subdomain, tunnelConfig.EndpointURL.Host)
	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		log.Debugf("Dial to local %s", localEndpoint.String())
		remote, errRemote := net.Dial("tcp", localEndpoint.String())
		client, errClient := listener.Accept()
		if errRemote == nil && errClient == nil && client != nil && remote != nil {
			go handleClient(client, remote)
		}
	}

}

func createTunnel(tunnelConfig *Config, semaphore *Semaphore) (net.Listener, error) {
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

	privateKey, err := readPrivateKeyFile(tunnelConfig.PrivateKeyPath)
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

	ncc, chans, reqs, err := ssh.NewClientConn(serverConn, serverEndpoint.String(), sshConfig)
	if err != nil {
		return listener, err
	}
	log.Debugf("SSH Connection Established via Jump %s -> %s", jumpServerEndpoint.String(), serverEndpoint.String())

	sClient := ssh.NewClient(ncc, chans, reqs)
	// Listen on remote server port
	listener, err = sClient.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open forwarding connection on remote server\n")
		return listener, err
	}
	log.Debugf("Open listen port on %s", remoteEndpoint.String())
	return listener, nil
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
func cleanup(config *Config) {
	fmt.Println("\nClosing tunnel")
	errSession := config.RestAPI.StartSession(config.RestAPI.RefreshToken)
	errDelete := config.RestAPI.DeleteTunnelAPI(config.Subdomain)
	if errSession != nil || errDelete != nil {
		fmt.Fprintf(os.Stderr,
			"We had some trouble deleting your tunnel. Use punch cleanup %s to make sure we know it's closed.\n", config.Subdomain)
	}

}
