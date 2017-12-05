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

	"k8s.io/kubectl/pkg/openapi"
)

func GetResources() (openapi.Resources, error) {
	schema, err := fakeSchema.OpenAPISchema()
	if err != nil {
		return nil, err
	}
	return openapi.NewOpenAPIData(schema)
}

func ExampleSchemaValidation_ValidateBytes() {
	resources, err := GetResources()
	if err != nil {
		return
	}
	validator := NewSchemaValidation(resources)

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

	err = validator.ValidateBytes(dataToValidate)
	if err != nil {
		return
	}
	fmt.Println("Valid data")
	// Output: Valid data
}
