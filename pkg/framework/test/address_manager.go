package test

import (
	"fmt"
	"net"
)

// AddressManager knows how to generate and remember a single address on some
// local interface for a service to listen on.
type AddressManager interface {
	Initialize() (port int, resolvedAddress string, err error)
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

// Initialize returns a address a process can listen on. It returns
// a tuple consisting of a free port and the hostname resolved to its IP.
func (d *DefaultAddressManager) Initialize() (port int, resolvedHost string, err error) {
	if d.port != 0 {
		return 0, "", fmt.Errorf("this DefaultAddressManager is already initialized")
	}
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
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
