package test

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

// EtcdConfig is a struct holding data to configure the Etcd process
type EtcdConfig struct {
	ClientURL string `valid:"required,url"`
	PeerURL   string `valid:"required,url"`
}

var etcdPortFinder = DefaultPortFinder

// NewEtcdConfig returns a simple config for Etcd with sane default values
func NewEtcdConfig() (etcdConfig *EtcdConfig, err error) {
	conf := &EtcdConfig{}
	host := "localhost"

	if port, addr, err := etcdPortFinder(host); err == nil {
		conf.ClientURL = fmt.Sprintf("http://%s:%d", addr, port)
	} else {
		return nil, err
	}

	if port, addr, err := etcdPortFinder(host); err == nil {
		conf.PeerURL = fmt.Sprintf("http://%s:%d", addr, port)
	} else {
		return nil, err
	}

	return conf, nil
}

// Validate checks that the config contains only valid URLs
func (c *EtcdConfig) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}
