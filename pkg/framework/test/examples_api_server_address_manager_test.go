package test_test

import (
	"fmt"

	. "k8s.io/kubectl/pkg/framework/test"
)

func ExampleAPIServer_firewalledAddressManager() {
	apiServer := &APIServer{
		AddressManager: &firewalledAddressManager{},
	}
	fmt.Println(apiServer)
}

type firewalledAddressManager struct {
	port int
	host string
}

func (a *firewalledAddressManager) Initialize() (int, string, error) {
	resolvedHost := resolveHostToIP("apiserver.integration.internal.local")
	if isAllowedIP(resolvedHost) {
		return 0, "", fmt.Errorf("the great firewall does not allow to bind services on this IP")
	}
	a.host = resolvedHost
	port, err := getFreePortFromAllowedRange()
	if err != nil {
		return 0, "", fmt.Errorf("could not find a free port on the range allowed by the great firewall")
	}
	a.port = port
	return a.port, a.host, nil
}
func (a *firewalledAddressManager) Port() (int, error) {
	return a.port, nil
}
func (a *firewalledAddressManager) Host() (string, error) {
	return a.host, nil
}

func resolveHostToIP(host string) string        { return "10.0.0.12" } // Look up host in local DNS
func isAllowedIP(ip string) bool                { return true }
func getFreePortFromAllowedRange() (int, error) { return 1234, nil }
