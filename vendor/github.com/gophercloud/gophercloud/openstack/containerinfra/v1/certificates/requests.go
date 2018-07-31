package certificates

import (
	"github.com/gophercloud/gophercloud"
)

// Get makes a request against the API to get details for a certificate.
func Get(client *gophercloud.ServiceClient, clusterID string) (r GetResult) {
	url := getURL(client, clusterID)

	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})

	return
}
