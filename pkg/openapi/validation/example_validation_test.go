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

package validation

import (
	"fmt"
	"path/filepath"

	tst "k8s.io/kubectl/pkg/openapi/openapitest"
)

func ExampleSchemaValidation_ValidateBytes() {
	validator := NewSchemaValidation(tst.NewFakeResources(filepath.Join("..", "openapitest", "swagger_test.json")))

	dataToValidate := []byte(`
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    name: redis-master
  name: name
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - image: redis
        name: redis
`)

	err := validator.ValidateBytes(dataToValidate)
	if err != nil {
		return
	}
	fmt.Println("Valid data")
	// Output: Valid data
}
