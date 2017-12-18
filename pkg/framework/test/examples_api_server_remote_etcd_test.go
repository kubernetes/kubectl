package test_test

import (
	"fmt"

	. "k8s.io/kubectl/pkg/framework/test"
)

func ExampleAPIServer_remoteEtcd() {
	apiServer := &APIServer{
		Etcd: &RemoteEtcd{},
	}
	fmt.Println(apiServer.Etcd.URL())
	// Output: https://my.big.internal.etcd.org:2379 <nil>
}

type RemoteEtcd struct{}

func (e *RemoteEtcd) URL() (string, error) {
	return "https://my.big.internal.etcd.org:2379", nil
}

func (e *RemoteEtcd) Start() error { return nil /* noop */ }
func (e *RemoteEtcd) Stop() error  { return nil /* noop */ }
