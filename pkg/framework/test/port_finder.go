package test

import (
	"fmt"
	"net"
)

// AddressManager knows how to generate and remember a single address on some
// local interface for a service to listen on.
type AddressManager interface {
	Initialize(host string) (port int, resolvedAddress string, err error)
	Host() (string, error)
	Port() (int, error)
}

//go:generate counterfeiter . AddressManager

// DefaultAddressManager implements an AddressManager. It allocates a new address
// (interface & port) a process can bind and keeps track of that.
type DefaultAddressManager struct {
	port int
	host string
}

// Initialize returns a address a process can listen on. Given a hostname it returns an address,
// a tuple consisting of a free port and the hostname resolved to its IP.
func (d *DefaultAddressManager) Initialize(host string) (port int, resolvedHost string, err error) {
	if d.port != 0 {
		return 0, "", fmt.Errorf("this DefaultAddressManager is already initialized")
	}
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", host))
	if err != nil {
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	d.port = l.Addr().(*net.TCPAddr).Port
	defer func() {
		err = l.Close()
	}()
	d.host = addr.IP.String()
	return d.port, d.host, nil
}

// Port returns the port that this DefaultAddressManager is managing. Port returns an
// error if this DefaultAddressManager has not yet been initialized.
func (d *DefaultAddressManager) Port() (int, error) {
	if d.port == 0 {
		return 0, fmt.Errorf("this DefaultAdressManager has is not initialized yet")
	}
	return d.port, nil
}

// Host returns the host that this DefaultAddressManager is managing. Host returns an
// error if this DefaultAddressManager has not yet been initialized.
func (d *DefaultAddressManager) Host() (string, error) {
	if d.host == "" {
		return "", fmt.Errorf("this DefaultAdressManager has is not initialized yet")
	}
	return d.host, nil
}

// PortFinder is the signature of a function returning a free port and a resolved
// address to listen/bind on, and an erro in case there were some problems finding
// a free pair of address & port
type PortFinder func(host string) (port int, resolvedAddress string, err error)

//go:generate counterfeiter . PortFinder

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
