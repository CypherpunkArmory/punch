package tunnel

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/tj/go-spin"
	"os"
	"os/signal"
	"strings"
	"time"
	//"strings"
	"sync"
	"syscall"
	//"time"
)

type outputConfig struct {
	connectionConfig *Config
	ConnectionUp     bool
}

func (o *outputConfig) String() string {
	var outputString string
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	if o.connectionConfig.EndpointType == "tcp" {
		outputString = fmt.Sprintf("\rAccess your service at %s://tcp.%s:%s",
			o.connectionConfig.EndpointType, o.connectionConfig.EndpointURL.Host, o.connectionConfig.TCPPorts[0])
	} else {
		outputString = fmt.Sprintf("\rAccess your website at %s://%s.%s",
			o.connectionConfig.EndpointType, o.connectionConfig.Subdomain, o.connectionConfig.EndpointURL.Host)
	}
	if o.connectionConfig.ConnectionUp {
		outputString = outputString + green(" •")
	} else {
		outputString = outputString + red(" •")
	}
	return outputString
}

func printConnections(printConfigs []outputConfig, connectionComplete *sync.WaitGroup) {
	commandRunning := true

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
		commandRunning = false
	}()
	connectionComplete.Wait()

	outputStrings := make([]string, len(printConfigs))
	screen.Clear()
	for commandRunning {
		for i, config := range printConfigs {
			outputStrings[i] = config.String()
		}
		screen.MoveTopLeft()
		fmt.Print(strings.Join(outputStrings[:], "\n"))
		time.Sleep(time.Millisecond * 500)
	}
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
