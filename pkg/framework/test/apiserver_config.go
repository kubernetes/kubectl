package test

import "github.com/asaskevich/govalidator"

// APIServerConfig is a struct holding data to configure the API Server process
type APIServerConfig struct {
	EtcdURL      string `valid:"required,url"`
	APIServerURL string `valid:"required,url"`
}

func (c *APIServerConfig) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}
