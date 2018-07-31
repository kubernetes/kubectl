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

package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/runtime"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmscheme "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1alpha3 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha3"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	"k8s.io/kubernetes/pkg/util/version"
)

// AnyConfigFileAndDefaultsToInternal reads either a InitConfiguration or JoinConfiguration and unmarshals it
func AnyConfigFileAndDefaultsToInternal(cfgPath string) (runtime.Object, error) {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	gvks, err := kubeadmutil.GroupVersionKindsFromBytes(b)
	if err != nil {
		return nil, err
	}

	// First, check if the gvk list has InitConfiguration and in that case try to unmarshal it
	if kubeadmutil.GroupVersionKindsHasInitConfiguration(gvks...) {
		return ConfigFileAndDefaultsToInternalConfig(cfgPath, &kubeadmapiv1alpha3.InitConfiguration{})
	}
	if kubeadmutil.GroupVersionKindsHasJoinConfiguration(gvks...) {
		return NodeConfigFileAndDefaultsToInternalConfig(cfgPath, &kubeadmapiv1alpha3.JoinConfiguration{})
	}
	return nil, fmt.Errorf("didn't recognize types with GroupVersionKind: %v", gvks)
}

// MarshalKubeadmConfigObject marshals an Object registered in the kubeadm scheme. If the object is a InitConfiguration, some extra logic is run
func MarshalKubeadmConfigObject(obj runtime.Object) ([]byte, error) {
	switch internalcfg := obj.(type) {
	case *kubeadmapi.InitConfiguration:
		return MarshalInitConfigurationToBytes(internalcfg, kubeadmapiv1alpha3.SchemeGroupVersion)
	default:
		return kubeadmutil.MarshalToYamlForCodecs(obj, kubeadmapiv1alpha3.SchemeGroupVersion, kubeadmscheme.Codecs)
	}
}

// DetectUnsupportedVersion reads YAML bytes, extracts the TypeMeta information and errors out with an user-friendly message if the API spec is too old for this kubeadm version
func DetectUnsupportedVersion(b []byte) error {
	gvks, err := kubeadmutil.GroupVersionKindsFromBytes(b)
	if err != nil {
		return err
	}

	// TODO: On our way to making the kubeadm API beta and higher, give good user output in case they use an old config file with a new kubeadm version, and
	// tell them how to upgrade. The support matrix will look something like this now and in the future:
	// v1.10 and earlier: v1alpha1
	// v1.11: v1alpha1 read-only, writes only v1alpha2 config
	// v1.12: v1alpha2 read-only, writes only v1beta1 config. Warns if the user tries to use v1alpha1
	// v1.13 and v1.14: v1beta1 read-only, writes only v1 config. Warns if the user tries to use v1alpha1 or v1alpha2.
	// v1.15: v1 is the only supported format.
	oldKnownAPIVersions := map[string]string{
		"kubeadm.k8s.io/v1alpha1": "v1.11",
	}
	// If we find an old API version in this gvk list, error out and tell the user why this doesn't work
	knownKinds := map[string]bool{}
	for _, gvk := range gvks {
		if useKubeadmVersion := oldKnownAPIVersions[gvk.GroupVersion().String()]; len(useKubeadmVersion) != 0 {
			return fmt.Errorf("your configuration file uses an old API spec: %q. Please use kubeadm %s instead and run 'kubeadm config migrate --old-config old.yaml --new-config new.yaml', which will write the new, similar spec using a newer API version.", gvk.GroupVersion().String(), useKubeadmVersion)
		}
		knownKinds[gvk.Kind] = true
	}
	// InitConfiguration, MasterConfiguration and NodeConfiguration are mutually exclusive, error if more than one are specified
	mutuallyExclusive := []string{constants.InitConfigurationKind, constants.MasterConfigurationKind, constants.JoinConfigurationKind, constants.NodeConfigurationKind}
	foundOne := false
	for _, kind := range mutuallyExclusive {
		if knownKinds[kind] {
			if !foundOne {
				foundOne = true
				continue
			}

			return fmt.Errorf("invalid configuration: kinds %v are mutually exclusive", mutuallyExclusive)
		}
	}
	return nil
}

// NormalizeKubernetesVersion resolves version labels, sets alternative
// image registry if requested for CI builds, and validates minimal
// version that kubeadm SetInitDynamicDefaultssupports.
func NormalizeKubernetesVersion(cfg *kubeadmapi.InitConfiguration) error {
	// Requested version is automatic CI build, thus use KubernetesCI Image Repository for core images
	if kubeadmutil.KubernetesIsCIVersion(cfg.KubernetesVersion) {
		cfg.CIImageRepository = constants.DefaultCIImageRepository
	}

	// Parse and validate the version argument and resolve possible CI version labels
	ver, err := kubeadmutil.KubernetesReleaseVersion(cfg.KubernetesVersion)
	if err != nil {
		return err
	}
	cfg.KubernetesVersion = ver

	// Parse the given kubernetes version and make sure it's higher than the lowest supported
	k8sVersion, err := version.ParseSemantic(cfg.KubernetesVersion)
	if err != nil {
		return fmt.Errorf("couldn't parse kubernetes version %q: %v", cfg.KubernetesVersion, err)
	}
	if k8sVersion.LessThan(constants.MinimumControlPlaneVersion) {
		return fmt.Errorf("this version of kubeadm only supports deploying clusters with the control plane version >= %s. Current version: %s", constants.MinimumControlPlaneVersion.String(), cfg.KubernetesVersion)
	}
	return nil
}

// LowercaseSANs can be used to force all SANs to be lowercase so it passes IsDNS1123Subdomain
func LowercaseSANs(sans []string) {
	for i, san := range sans {
		lowercase := strings.ToLower(san)
		if lowercase != san {
			glog.V(1).Infof("lowercasing SAN %q to %q", san, lowercase)
			sans[i] = lowercase
		}
	}
}
