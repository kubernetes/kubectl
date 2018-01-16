// Package test an integration test framework for k8s
package test

import (
	"net/url"
)

// ControlPlane is a struct that knows how to start your test control plane.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in future.
type ControlPlane struct {
	APIServer *APIServer
	Etcd      *Etcd
}

// NewControlPlane will give you a ControlPlane struct that's properly wired together.
func NewControlPlane() *ControlPlane {
	etcd := &Etcd{}
	apiserver := &APIServer{}

	return &ControlPlane{
		APIServer: apiserver,
		Etcd:      etcd,
	}
}

// Start will start your control plane. To stop it, call Stop().
func (f *ControlPlane) Start() error {
	if err := f.Etcd.Start(); err != nil {
		return err
	}

	f.APIServer.EtcdURL = &f.Etcd.processState.URL
	return f.APIServer.Start()
}

// Stop will stop your control plane, and clean up their data.
func (f *ControlPlane) Stop() error {
	if err := f.APIServer.Stop(); err != nil {
		return err
	}
	return f.Etcd.Stop()
}

// APIURL returns the URL you should connect to to talk to your API.
func (f *ControlPlane) APIURL() url.URL {
	return f.APIServer.processState.URL
}
