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

package kinflate

import (
	"os"
	"path"
	"path/filepath"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
)

func adjustPathsForConfigMapAndSecret(in *manifest.Manifest, pathToDir []string) *resource {
	out := &resource{}
	out.configmaps = adjustPathForConfigMaps(in.Configmaps, pathToDir)
	out.secrets = adjustPathForSecrets(in.Secrets, pathToDir)
	return out
}

func adjustPathForConfigMaps(cms []manifest.ConfigMap, prefix []string) []manifest.ConfigMap {
	for i, cm := range cms {
		if len(cm.FileSources) > 0 {
			for j, fileSource := range cm.FileSources {
				cms[i].FileSources[j] = adjustPath(fileSource, prefix)
			}
		}
		if len(cm.EnvSource) > 0 {
			cms[i].EnvSource = adjustPath(cm.EnvSource, prefix)
		}
	}
	return cms
}

func adjustPathForSecrets(secrets []manifest.Secret, prefix []string) []manifest.Secret {
	for i, secret := range secrets {
		if len(secret.FileSources) > 0 {
			for j, fileSource := range secret.FileSources {
				secrets[i].FileSources[j] = adjustPath(fileSource, prefix)
			}
		}
		if len(secret.EnvSource) > 0 {
			secrets[i].EnvSource = adjustPath(secret.EnvSource, prefix)
		}
		if secret.TLS != nil {
			secrets[i].TLS.CertFile = adjustPath(secret.TLS.CertFile, prefix)
			secrets[i].TLS.KeyFile = adjustPath(secret.TLS.KeyFile, prefix)
		}
	}
	return secrets
}

func adjustPath(original string, prefix []string) string {
	return path.Join(append(prefix, original)...)
}

func adjustPaths(original []string, prefix []string) ([]string, error) {
	adjusted := []string{}
	var e error
	for _, filename := range original {
		filepath.Walk(adjustPath(filename, prefix), func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				e = err
				return err
			}
			adjusted = append(adjusted, path)
			return nil
		})
	}
	return adjusted, e
}
