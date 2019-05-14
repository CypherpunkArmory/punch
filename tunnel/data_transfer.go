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
	"io"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	ioFinished := &sync.WaitGroup{}
	ioFinished.Add(2)
	errorCh := make(chan error, 2)

	go copyData(client, "client", remote, "remote", ioFinished, errorCh)
	go copyData(remote, "remote", client, "client", ioFinished, errorCh)
	ioFinished.Wait()

	select {
	case firstError := <-errorCh:
		log.Debugf(firstError.Error())
	default:
		return
	}
}

func copyData(dst net.Conn, dstName string, src net.Conn, srcName string, done *sync.WaitGroup, errorCh chan error) {
	defer done.Done()
	amt, err := io.Copy(dst, src)
	if err != nil {
		errorCh <- fmt.Errorf(
			"%s -> %s error: %s",
			srcName,
			dstName,
			err.Error())
	}
	log.Debugf("%s <- %s (%d bytes)", dstName, srcName, amt)
	src.Close()
	dst.Close()
}
