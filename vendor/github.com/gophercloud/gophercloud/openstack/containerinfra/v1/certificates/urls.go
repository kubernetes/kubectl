package certificates

import (
	"github.com/gophercloud/gophercloud"
)

var apiName = "certificates"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiName)
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL(apiName, id)
}
