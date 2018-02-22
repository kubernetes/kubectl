package resource

import (
	"testing"

	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/loader"
)

var encoded = []byte(`apiVersion: v1
kind: Deployment
metadata:
  name: dply1
---
apiVersion: v1
kind: Deployment
metadata:
  name: dply2
`)

type fakeLoader struct {
}

func (l fakeLoader) New(newRoot string) (loader.Loader, error) {
	return l, nil
}

func (l fakeLoader) Load(location string) ([]byte, error) {
	return encoded, nil
}

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

func TestAppResourceList_Resources(t *testing.T) {
	l := fakeLoader{}
	expected := []*Resource{
		{Data: makeUnconstructed("dply1")},
		{Data: makeUnconstructed("dply2")},
	}

	resources, _ := ResourcesFromPath("fake/path", l)
	if len(resources) != 2 {
		t.Fatalf("%#v should contain 2 appResource, but got %d", resources, len(resources))
	}

	for i, r := range resources {
		if !reflect.DeepEqual(r.Data, expected[i].Data) {
			t.Fatalf("expected %v, but got %v", expected[i].Data, r.Data)
		}
	}
}
