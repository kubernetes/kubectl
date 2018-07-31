/*
Copyright 2016 The Kubernetes Authors.

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

package images

import (
	"fmt"
	"runtime"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
)

// GetGenericImage generates and returns a platform agnostic image (backed by manifest list)
func GetGenericImage(prefix, image, tag string) string {
	return fmt.Sprintf("%s/%s:%s", prefix, image, tag)
}

// GetGenericArchImage generates and returns an image based on the current runtime arch
func GetGenericArchImage(prefix, image, tag string) string {
	return fmt.Sprintf("%s/%s-%s:%s", prefix, image, runtime.GOARCH, tag)
}

// GetKubeControlPlaneImageNoOverride generates and returns the image for the core Kubernetes components ignoring the unified control plane image
func GetKubeControlPlaneImageNoOverride(image string, cfg *kubeadmapi.InitConfiguration) string {
	repoPrefix := cfg.GetControlPlaneImageRepository()
	kubernetesImageTag := kubeadmutil.KubernetesVersionToImageTag(cfg.KubernetesVersion)
	return GetGenericArchImage(repoPrefix, image, kubernetesImageTag)
}

// GetKubeControlPlaneImage generates and returns the image for the core Kubernetes components or returns the unified control plane image if specified
func GetKubeControlPlaneImage(image string, cfg *kubeadmapi.InitConfiguration) string {
	if cfg.UnifiedControlPlaneImage != "" {
		return cfg.UnifiedControlPlaneImage
	}
	return GetKubeControlPlaneImageNoOverride(image, cfg)
}

// GetEtcdImage generates and returns the image for etcd or returns cfg.Etcd.Local.Image if specified
func GetEtcdImage(cfg *kubeadmapi.InitConfiguration) string {
	if cfg.Etcd.Local != nil && cfg.Etcd.Local.Image != "" {
		return cfg.Etcd.Local.Image
	}
	etcdImageTag := constants.DefaultEtcdVersion
	etcdImageVersion, err := constants.EtcdSupportedVersion(cfg.KubernetesVersion)
	if err == nil {
		etcdImageTag = etcdImageVersion.String()
	}
	return GetGenericArchImage(cfg.ImageRepository, constants.Etcd, etcdImageTag)
}

// GetAllImages returns a list of container images kubeadm expects to use on a control plane node
func GetAllImages(cfg *kubeadmapi.InitConfiguration) []string {
	imgs := []string{}
	imgs = append(imgs, GetKubeControlPlaneImage(constants.KubeAPIServer, cfg))
	imgs = append(imgs, GetKubeControlPlaneImage(constants.KubeControllerManager, cfg))
	imgs = append(imgs, GetKubeControlPlaneImage(constants.KubeScheduler, cfg))
	imgs = append(imgs, GetKubeControlPlaneImageNoOverride(constants.KubeProxy, cfg))

	// pause, etcd and kube-dns are not available on the ci image repository so use the default image repository.
	imgs = append(imgs, GetGenericImage(cfg.ImageRepository, "pause", "3.1"))

	// if etcd is not external then add the image as it will be required
	if cfg.Etcd.Local != nil {
		imgs = append(imgs, GetEtcdImage(cfg))
	}

	dnsImage := GetGenericArchImage(cfg.ImageRepository, "k8s-dns-kube-dns", constants.KubeDNSVersion)
	if features.Enabled(cfg.FeatureGates, features.CoreDNS) {
		dnsImage = fmt.Sprintf("%s/%s:%s", cfg.ImageRepository, constants.CoreDNS, constants.CoreDNSVersion)
	}
	imgs = append(imgs, dnsImage)
	return imgs
}
