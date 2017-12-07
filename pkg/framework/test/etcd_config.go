package test

import "github.com/asaskevich/govalidator"

// EtcdConfig is a struct holding data to configure the Etcd process
type EtcdConfig struct {
	ClientURL string `valid:"required,url"`
	PeerURL   string `valid:"required,url"`
}

// Validate checks that the config contains only valid URLs
func (c *EtcdConfig) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}
