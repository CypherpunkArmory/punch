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
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

//StartReverseTunnel Main tunneling function. Handles connections and forwarding
func StartReverseTunnel(tunnelConfig ...Config) {
	fmt.Println("Use Ctrl-c to close the tunnels")
	semaphore := Semaphore{}
	wg := sync.WaitGroup{}
	wg.Add(len(tunnelConfig))
	jumpConn, err := connectToJumpHost(&tunnelConfig[0], &semaphore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "We could not connect to the jump host")
		cleanup(&tunnelConfig[0])
	}
	for i := range tunnelConfig {
		go startReverseTunnel(jumpConn, &tunnelConfig[i], &wg, &semaphore, tunnelConfig[i].TCPPorts[i])
	}
	wg.Wait()
}

func startReverseTunnel(jumpConn *ssh.Client, tunnelConfig *Config, wg *sync.WaitGroup, semaphore *Semaphore, tcpPort string) {
	defer cleanup(tunnelConfig)

	if wg != nil {
		defer wg.Done()
	}
	sClient, err := createTunnel(jumpConn, tunnelConfig, semaphore)
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
		<-startCloseChannel
		if semaphore.CanRun() {
			cleanup(tunnelConfig)
			os.Exit(0)
		}
	}()
	var outputString string
	if tunnelConfig.EndpointType == "tcp" {
		outputString = fmt.Sprintf("Access your service at %s://tcp.%s:%s",
			tunnelConfig.EndpointType, tunnelConfig.EndpointURL.Host, tcpPort)
	} else {
		outputString = fmt.Sprintf("Access your website at %s://%s.%s",
			tunnelConfig.EndpointType, tunnelConfig.Subdomain, tunnelConfig.EndpointURL.Host)
	}
	fmt.Println(outputString)
	var listener net.Listener
	listener, err = sClient.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		log.Debugf("Dial to local %s", localEndpoint.String())
		_, errInitialRemoteConnect := net.Dial("tcp", localEndpoint.String())
		if errInitialRemoteConnect == nil {
			client, errClient := listener.Accept()
			remote, errRemote := net.Dial("tcp", localEndpoint.String())
			if errRemote == nil && errClient == nil && client != nil && remote != nil {
				// start goroutine
				go handleClient(client, remote)
			}
		} else {
			log.Debugf("Err %s", errInitialRemoteConnect.Error())
			listener.Close()
			listener = nil // you can't close the underlying file descriptor on the connection
			// so you need to let the listener be GC'ed by replacing it with a new object
			log.Debugf("No local listener")
			time.Sleep(1000 * time.Millisecond)
			log.Debugf("Trying again")
			listener, err = sClient.Listen("tcp", remoteEndpoint.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error connecting to local host")
				startCloseChannel <- syscall.SIGINT
			}
		}

	}

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
