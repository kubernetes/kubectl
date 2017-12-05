/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openapi_test

import (
	"fmt"

	"k8s.io/kubectl/pkg/openapi"
	tst "k8s.io/kubectl/pkg/openapi/testing"
)

func ExampleGetter_Get() {
	client := tst.NewFakeClient(&fakeSchema)

	// Pass client to obtain getter
	getter := openapi.NewOpenAPIGetter(client)

	// Get saved recources
	resources, err := getter.Get()
	if err == nil && resources != nil {
		fmt.Println("Valid resources")
	}
	// Output: Valid resources
}
