package tunnel

import (
	"fmt"
	"io"
	"sync"
)

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client io.ReadWriteCloser, remote io.ReadWriteCloser) {
	ioFinished := &sync.WaitGroup{}
	ioFinished.Add(2)
	errorCh := make(chan error, 2)

	go copyData(client, "client", remote, "remote", ioFinished, errorCh)
	go copyData(remote, "remote", client, "client", ioFinished, errorCh)

	ioFinished.Wait()
	err := <-errorCh
	if err != nil {
	}
}

func copyData(dst io.WriteCloser, dstName string, src io.ReadCloser, srcName string, done *sync.WaitGroup, errorCh chan error) {
	defer done.Done()

	if _, err := io.Copy(dst, src); err != nil {
		errorCh <- fmt.Errorf(
			"%s -> %s error: %s",
			srcName,
			dstName,
			err.Error())
	}
	src.Close()
	dst.Close()
}
