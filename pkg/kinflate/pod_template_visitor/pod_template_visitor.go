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

package pod_template_visitor

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kapps "k8s.io/kubectl/pkg/kinflate/apps"
)

type PodTemplateSpecVisitor struct {
	Object  *unstructured.Unstructured
	MungeFn func(podTemplateSpec map[string]interface{}) error
	Err     error
}

var _ kapps.KindVisitor = &PodTemplateSpecVisitor{}

func (v *PodTemplateSpecVisitor) VisitDeployment(kind kapps.GroupKindElement) {
	v.Err = v.mungePodTemplateSpec([]string{"spec", "template"})
}

func (v *PodTemplateSpecVisitor) VisitStatefulSet(kind kapps.GroupKindElement) {
	v.Err = v.mungePodTemplateSpec([]string{"spec", "template"})
}

func (v *PodTemplateSpecVisitor) VisitDaemonSet(kind kapps.GroupKindElement) {
	v.Err = v.mungePodTemplateSpec([]string{"spec", "template"})
}

func (v *PodTemplateSpecVisitor) VisitJob(kind kapps.GroupKindElement) {
	v.Err = v.mungePodTemplateSpec([]string{"spec", "template"})
}

func (v *PodTemplateSpecVisitor) VisitReplicaSet(kind kapps.GroupKindElement) {
	v.Err = v.mungePodTemplateSpec([]string{"spec", "template"})
}

func (v *PodTemplateSpecVisitor) VisitPod(kind kapps.GroupKindElement) {}

func (v *PodTemplateSpecVisitor) VisitReplicationController(kind kapps.GroupKindElement) {
	v.Err = v.mungePodTemplateSpec([]string{"spec", "template"})
}

func (v *PodTemplateSpecVisitor) VisitCronJob(kind kapps.GroupKindElement) {
	v.Err = v.mungePodTemplateSpec([]string{"spec", "jobTemplate", "spec", "template"})
}

func walkMapPath(start map[string]interface{}, path []string) (map[string]interface{}, error) {
	finish := start
	for i := 0; i < len(path); i++ {
		var ok bool
		finish, ok = finish[path[i]].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("key:%s of path:%v not found in map:%v", path[i], path, start)
		}
	}

	return finish, nil
}

func (v *PodTemplateSpecVisitor) mungePodTemplateSpec(pathToPodTemplateSpec []string) error {
	obj := v.Object.UnstructuredContent()
	podTemplateSpec, err := walkMapPath(obj, pathToPodTemplateSpec)
	if err != nil {
		return err
	}
	return v.MungeFn(podTemplateSpec)
}
