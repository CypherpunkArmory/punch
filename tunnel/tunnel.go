// Punch CLI used for interacting with holepunch.io
// Copyright (C) 2018-2019  Orb.House, LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package tunnel

import (
	"errors"
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

func internalEndpoint(endpointType string) (*Endpoint, error) {
	if endpointType == "http" {
		return &Endpoint{
			Host: "localhost", // localhost here is the remote SSHD daemon container
			Port: "3000",
		}, nil
	} else if endpointType == "https" {
		return &Endpoint{
			Host: "localhost", // localhost here is the remote SSHD daemon container
			Port: "3001",
		}, nil
	} else {
		return nil, errors.New("unknown Endpoint Type")
	}
}

//StartReverseTunnel Main tunneling function. Handles connections and forwarding
func StartReverseTunnel(tunnelConfig *Config, wg *sync.WaitGroup, semaphore *Semaphore) {
	defer cleanup(tunnelConfig)

	if wg != nil {
		defer wg.Done()
	}
	sClient, err := createTunnel(tunnelConfig, semaphore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	remoteEndpoint, err := internalEndpoint(tunnelConfig.EndpointType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	defer sClient.Close()

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
		log.Debugf("Closing Channel")
		<-startCloseChannel
		if semaphore.CanRun() {
			cleanup(tunnelConfig)
			os.Exit(0)
		}
	}()

	fmt.Printf("Access your website at %s://%s.%s\n",
		tunnelConfig.EndpointType, tunnelConfig.Subdomain, tunnelConfig.EndpointURL.Host)
	// handle incoming connections on reverse forwarded tunnel
	var listener net.Listener
	listener, err = sClient.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		log.Debugf("Dial to local %s", localEndpoint.String())
		log.Debugf("%+v", listener)
		_, errInitialRemoteConnect := net.Dial("tcp", localEndpoint.String())
		if errInitialRemoteConnect == nil {
			fmt.Printf("connected\n")
			client, errClient := listener.Accept()
			remote, errRemote := net.Dial("tcp", localEndpoint.String())
			if errRemote == nil && errClient == nil && client != nil && remote != nil {
				// start goroutine
				go handleClient(client, remote)
			}
		} else {
			fmt.Printf("dis-connected\n")
			log.Debugf("Err %s", errInitialRemoteConnect.Error())
			listener.Close()
			listener = nil // you can't close the underlying file descriptor on the connection
			// so you need to let the listener be GC'ed by replacing it with a new object
			log.Debugf("No local listener")
			time.Sleep(1000 * time.Millisecond)
			log.Debugf("Trying again")
			listener, _ = sClient.Listen("tcp", remoteEndpoint.String())
		}

	}

}

func createTunnel(tunnelConfig *Config, semaphore *Semaphore) (*ssh.Client, error) {
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
	// FIXME:  This should be a LUT

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

	privateKey, err := readPrivateKeyFile(tunnelConfig.PrivateKeyPath)
	if err != nil {
		return nil, err
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
	log.Debugf("SSH Connection Established via Jump %s -> %s", jumpServerEndpoint.String(), serverEndpoint.String())

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
func cleanup(config *Config) {
	fmt.Println("\nClosing tunnel")
	errSession := config.RestAPI.StartSession(config.RestAPI.RefreshToken)
	errDelete := config.RestAPI.DeleteTunnelAPI(config.Subdomain)
	if errSession != nil || errDelete != nil {
		fmt.Fprintf(os.Stderr,
			"We had some trouble deleting your tunnel. Use punch cleanup %s to make sure we know it's closed.\n", config.Subdomain)
	}
}
