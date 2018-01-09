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

package kinflate

import (
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
)

type testCase struct {
	description string
	in          []byte
	expected    []byte
	pathToField []string
	fn          func(interface{}) (interface{}, error)
}

var input = []byte(`
metadata:
  labels:
    app: foo
spec:
  containers:
  - name: nginx
    image: nginx:1.7.9
  - name: busybox
    image: busybox
    volumeMounts:
    - mountPath: /tmp/env
      name: app-env
    - mountPath: /tmp/tls
      name: app-tls
  selector:
    app: foo
  volumes:
  - configMap:
      name: app-env
    name: app-env
  - secret:
      name: app-tls
    name: app-tls
`)

var testCases = []testCase{
	{
		description: "add prefix",
		in:          input,
		expected: []byte(`
metadata:
  labels:
    app: foo
spec:
  containers:
  - image: nginx:1.7.9
    name: nginx
  - image: busybox
    name: busybox
    volumeMounts:
    - mountPath: /tmp/env
      name: app-env
    - mountPath: /tmp/tls
      name: app-tls
  selector:
    app: foo
  volumes:
  - configMap:
      name: someprefix-app-env
    name: app-env
  - name: app-tls
    secret:
      name: app-tls
`),
		pathToField: []string{"spec", "volumes", "configMap", "name"},
		fn:          addPrefix("someprefix-"),
	},
	{
		description: "add map",
		in:          input,
		expected: []byte(`
metadata:
  labels:
    app: foo
spec:
  containers:
  - image: nginx:1.7.9
    name: nginx
  - image: busybox
    name: busybox
    volumeMounts:
    - mountPath: /tmp/env
      name: app-env
    - mountPath: /tmp/tls
      name: app-tls
  selector:
    app: foo
    component: bar
  volumes:
  - configMap:
      name: app-env
    name: app-env
  - name: app-tls
    secret:
      name: app-tls
`),
		pathToField: []string{"spec", "selector"},
		fn:          addMap(map[string]string{"component": "bar"}),
	},
}

func TestUpdate(t *testing.T) {
	for _, tc := range testCases {
		var m, expected map[string]interface{}
		err := yaml.Unmarshal(tc.in, &m)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err = mutateField(m, tc.pathToField, tc.fn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err = yaml.Unmarshal(tc.expected, &expected)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(m, expected) {
			t.Fatalf("in testcase: %q updated:\n%s\ndoesn't match expected:\n%s\n", tc.description, m, expected)
		}
	}

}
