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

package commands

import (
	"fmt"
	"io"

	"k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/kinflate/util/fs"

	"github.com/spf13/cobra"
)

func NewCmdAddConfigMap(errOut io.Writer, fsys fs.FileSystem) *cobra.Command {
	var config dataConfig
	cmd := &cobra.Command{
		Use:   "configmap NAME [--from-file=[key=]source] [--from-literal=key1=value1]",
		Short: "Adds a configmap to your manifest file.",
		Long:  "",
		Example: `
	# Adds a configmap to the Manifest (with a specified key)
	kinflate configmap my-configmap --from-file=my-key=file/path --from-literal=my-literal=12345

	# Adds a configmap to the Manifest (key is the filename)
	kinflate configmap my-configmap --from-file=file/path

	# Adds a configmap from env-file
	kinflate configmap my-configmap --from-env-file=env/path.env
`,
		RunE: func(_ *cobra.Command, args []string) error {
			err := config.Validate(args)
			if err != nil {
				return err
			}

			// TODO(apelisse,droot): Do something with that config.

			loader := kutil.ManifestLoader{FS: fsys}

			m, err := loader.Read(kinflate.KubeManifestFileName)
			if err != nil {
				return err
			}

			err = addConfigMap(m, config)
			if err != nil {
				return err
			}

			return loader.Write(kinflate.KubeManifestFileName, m)

		},
	}

	cmd.Flags().StringSliceVar(&config.FileSources, "from-file", []string{}, "Key file can be specified using its file path, in which case file basename will be used as configmap key, or optionally with a key and file path, in which case the given key will be used.  Specifying a directory will iterate each named file in the directory whose basename is a valid configmap key.")
	cmd.Flags().StringArrayVar(&config.LiteralSources, "from-literal", []string{}, "Specify a key and literal value to insert in configmap (i.e. mykey=somevalue)")
	cmd.Flags().StringVar(&config.EnvFileSource, "from-env-file", "", "Specify the path to a file to read lines of key=val pairs to create a configmap (i.e. a Docker .env file).")

	return cmd
}

func addConfigMap(m *manifest.Manifest, config dataConfig) error {
	cm := locateConfigMap(m, config.Name)
	if cm == nil {
		cm = &v1alpha1.ConfigMap{}
	}

	// deep copy the configmap

	mergedContent, err := mergeData(&cm.Generic, config)
	if err != nil {
		return fmt.Errorf("unable to add configmap: %v", err)
	}

	cm.Generic = *mergedContent

	return nil
}

func locateConfigMap(m *manifest.Manifest, namePrefix string) *manifest.ConfigMap {
	for i, v := range m.Configmaps {
		if namePrefix == v.NamePrefix {
			return &m.Configmaps[i]
		}
	}
	return nil
}

func mergeData(src *manifest.Generic, config dataConfig) (*manifest.Generic, error) {
	return nil, nil
}
