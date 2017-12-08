package test

import (
	"fmt"
	"net"
)

// PortFinder is the signature of a function returning a free port and a resolved
// address to listen/bind on, and an erro in case there were some problems finding
// a free pair of address & port
type PortFinder func(host string) (port int, resolvedAddress string, err error)

// DefaultPortFinder is the default implementation of PortFinder: It resolves that
// hostname or IP handed in and asks the kernel for a random port. To make that a bit
// safer, it also tries to bind to that port to make sure it is not in use. If all goes
// well the port and the resolved address are returned, otherwise the error is forwarded.
func DefaultPortFinder(host string) (port int, resolvedAddress string, err error) {
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", host))
	if err != nil {
		return
	}
	resolvedAddress = addr.IP.String()
	var l *net.TCPListener
	l, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer func() {
		err = l.Close()
	}()

	return l.Addr().(*net.TCPAddr).Port, resolvedAddress, nil
}
