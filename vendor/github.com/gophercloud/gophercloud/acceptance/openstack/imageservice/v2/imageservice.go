// Package v2 contains common functions for creating imageservice resources
// for use in acceptance tests. See the `*_test.go` files for example usages.
package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	th "github.com/gophercloud/gophercloud/testhelper"
)

// CreateEmptyImage will create an image, but with no actual image data.
// An error will be returned if an image was unable to be created.
func CreateEmptyImage(t *testing.T, client *gophercloud.ServiceClient) (*images.Image, error) {
	var image *images.Image

	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create image: %s", name)

	protected := false
	visibility := images.ImageVisibilityPrivate
	createOpts := &images.CreateOpts{
		Name:            name,
		ContainerFormat: "bare",
		DiskFormat:      "qcow2",
		MinDisk:         0,
		MinRAM:          0,
		Protected:       &protected,
		Visibility:      &visibility,
		Properties: map[string]string{
			"architecture": "x86_64",
		},
		Tags: []string{"foo", "bar", "baz"},
	}

	image, err := images.Create(client, createOpts).Extract()
	if err != nil {
		return image, err
	}

	newImage, err := images.Get(client, image.ID).Extract()
	if err != nil {
		return image, err
	}

	t.Logf("Created image %s: %#v", name, newImage)

	th.CheckEquals(t, newImage.Name, name)
	th.CheckEquals(t, newImage.Properties["architecture"], "x86_64")
	return newImage, nil
}

// DeleteImage deletes an image.
// A fatal error will occur if the image failed to delete. This works best when
// used as a deferred function.
func DeleteImage(t *testing.T, client *gophercloud.ServiceClient, image *images.Image) {
	err := images.Delete(client, image.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete image %s: %v", image.ID, err)
	}

	t.Logf("Deleted image: %s", image.ID)
}
