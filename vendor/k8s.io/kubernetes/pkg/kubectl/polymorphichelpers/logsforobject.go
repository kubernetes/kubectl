/*
Copyright 2018 The Kubernetes Authors.

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

package polymorphichelpers

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	coreinternal "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions"
)

func logsForObject(restClientGetter genericclioptions.RESTClientGetter, object, options runtime.Object, timeout time.Duration, allContainers bool) ([]*rest.Request, error) {
	clientConfig, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := internalclientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	return logsForObjectWithClient(clientset, object, options, timeout, allContainers)
}

// this is split for easy test-ability
func logsForObjectWithClient(clientset internalclientset.Interface, object, options runtime.Object, timeout time.Duration, allContainers bool) ([]*rest.Request, error) {
	opts, ok := options.(*coreinternal.PodLogOptions)
	if !ok {
		return nil, errors.New("provided options object is not a PodLogOptions")
	}

	switch t := object.(type) {
	case *coreinternal.PodList:
		ret := []*rest.Request{}
		for i := range t.Items {
			currRet, err := logsForObjectWithClient(clientset, &t.Items[i], options, timeout, allContainers)
			if err != nil {
				return nil, err
			}
			ret = append(ret, currRet...)
		}
		return ret, nil

	case *corev1.PodList:
		ret := []*rest.Request{}
		for i := range t.Items {
			currRet, err := logsForObjectWithClient(clientset, &t.Items[i], options, timeout, allContainers)
			if err != nil {
				return nil, err
			}
			ret = append(ret, currRet...)
		}
		return ret, nil

	case *coreinternal.Pod:
		// if allContainers is true, then we're going to locate all containers and then iterate through them. At that point, "allContainers" is false
		if !allContainers {
			return []*rest.Request{clientset.Core().Pods(t.Namespace).GetLogs(t.Name, opts)}, nil
		}

		ret := []*rest.Request{}
		for _, c := range t.Spec.InitContainers {
			currOpts := opts.DeepCopy()
			currOpts.Container = c.Name
			currRet, err := logsForObjectWithClient(clientset, t, options, timeout, false)
			if err != nil {
				return nil, err
			}
			ret = append(ret, currRet...)
		}
		for _, c := range t.Spec.Containers {
			currOpts := opts.DeepCopy()
			currOpts.Container = c.Name
			currRet, err := logsForObjectWithClient(clientset, t, options, timeout, false)
			if err != nil {
				return nil, err
			}
			ret = append(ret, currRet...)
		}

		return ret, nil

	case *corev1.Pod:
		// if allContainers is true, then we're going to locate all containers and then iterate through them. At that point, "allContainers" is false
		if !allContainers {
			return []*rest.Request{clientset.Core().Pods(t.Namespace).GetLogs(t.Name, opts)}, nil
		}

		ret := []*rest.Request{}
		for _, c := range t.Spec.InitContainers {
			currOpts := opts.DeepCopy()
			currOpts.Container = c.Name
			currRet, err := logsForObjectWithClient(clientset, t, options, timeout, false)
			if err != nil {
				return nil, err
			}
			ret = append(ret, currRet...)
		}
		for _, c := range t.Spec.Containers {
			currOpts := opts.DeepCopy()
			currOpts.Container = c.Name
			currRet, err := logsForObjectWithClient(clientset, t, options, timeout, false)
			if err != nil {
				return nil, err
			}
			ret = append(ret, currRet...)
		}

		return ret, nil
	}

	namespace, selector, err := SelectorsForObject(object)
	if err != nil {
		return nil, fmt.Errorf("cannot get the logs from %T: %v", object, err)
	}
	sortBy := func(pods []*v1.Pod) sort.Interface { return controller.ByLogging(pods) }
	pod, numPods, err := GetFirstPod(clientset.Core(), namespace, selector.String(), timeout, sortBy)
	if err != nil {
		return nil, err
	}
	if numPods > 1 {
		fmt.Fprintf(os.Stderr, "Found %v pods, using pod/%v\n", numPods, pod.Name)
	}

	return logsForObjectWithClient(clientset, pod, options, timeout, allContainers)
}
