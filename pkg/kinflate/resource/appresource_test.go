package resource

import (
	"testing"

	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/loader/loadertest"
)

func makeUnconstructed(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

func TestNewFromPath(t *testing.T) {

	resourceStr := `apiVersion: v1
kind: Deployment
metadata:
  name: dply1
---
apiVersion: v1
kind: Deployment
metadata:
  name: dply2
`

	l := loadertest.NewFakeLoader("/home/seans/project")
	if ferr := l.AddFile("/home/seans/project/deployment.yaml", []byte(resourceStr)); ferr != nil {
		t.Fatalf("Error adding fake file: %v\n", ferr)
	}
	expected := []*Resource{
		{Data: makeUnconstructed("dply1")},
		{Data: makeUnconstructed("dply2")},
	}

	resources, _ := NewFromPath("/home/seans/project/deployment.yaml", l)
	if len(resources) != 2 {
		t.Fatalf("%#v should contain 2 appResource, but got %d", resources, len(resources))
	}

	for i, r := range resources {
		if !reflect.DeepEqual(r.Data, expected[i].Data) {
			t.Fatalf("expected %v, but got %v", expected[i].Data, r.Data)
		}
	}
}
