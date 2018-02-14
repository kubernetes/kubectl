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

package commands

import (
	"fmt"
	"io"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/configmapandsecret"
	"k8s.io/kubectl/pkg/kinflate/constants"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/kinflate/util/fs"

	"github.com/spf13/cobra"
)

// newCmdAddSecret returns a new Cobra command that wraps generic and tls secrets.
func newCmdAddSecret(errOut io.Writer, fsys fs.FileSystem) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Adds a secret using specified subcommand",
		Example: `
	# Adds a TLS secret to the Manifest (with a specified key)
	kinflate add secret tls my-tls-secret --cert=cert/path.cert --key=key/path.key
`,
	}
	cmd.AddCommand(
		newCmdAddSecretTLS(errOut, fsys),
	)

	return cmd
}

type addTLSSecret struct {
	// Name of secret (required)
	Name string
	// Cert is the file path to the cerificate (required)
	Cert string
	// Key is the file path to the key (required)
	Key string
}

// validate validates required fields are set to support structured generation.
func (a *addTLSSecret) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("name must be specified once")
	}
	a.Name = args[0]
	if len(a.Cert) == 0 {
		return fmt.Errorf("cert is required")
	}
	if len(a.Key) == 0 {
		return fmt.Errorf("key is required")
	}
	// TODO: Should we check if the path exists? if it's valid, if it's within the same (sub-)directory?
	return nil
}

// newCmdCreateSecretTLS is a macro command for creating secrets to work with Docker registries
func newCmdAddSecretTLS(errOut io.Writer, fsys fs.FileSystem) *cobra.Command {
	var config addTLSSecret
	cmd := &cobra.Command{
		Use:   "tls NAME --cert=path/to/cert/file --key=path/to/key/file",
		Short: "Adds a TLS secret.",
		Long:  "",
		Example: `
	# Adds a TLS secret to the Manifest (with a specified key)
	kinflate secret tls my-tls-secret --cert=cert/path.cert --key=key/path.key
`,
		RunE: func(_ *cobra.Command, args []string) error {
			err := config.Validate(args)
			if err != nil {
				return err
			}

			loader := kutil.ManifestLoader{FS: fsys}
			m, err := loader.Read(constants.KubeManifestFileName)
			if err != nil {
				return err
			}

			err = addTLSSecretToManifest(m, config)
			if err != nil {
				return err
			}

			return loader.Write(constants.KubeManifestFileName, m)
		},
	}

	cmd.Flags().StringVar(&config.Cert, "cert", "", "Path to PEM encoded public key certificate.")
	cmd.Flags().StringVar(&config.Key, "key", "", "Path to private key associated with given certificate.")

	return cmd
}

// addTLSSecretToManifest appends the TLS secret to the manifest, or returns
// an error if the secret already exists.
func addTLSSecretToManifest(m *manifest.Manifest, a addTLSSecret) error {

	if tlsSecretExists(m, a.Name) {
		return fmt.Errorf("TLS Secret already exists")
	}

	tls := manifest.TLSSecret{
		Name:     a.Name,
		CertFile: a.Cert,
		KeyFile:  a.Key,
	}
	m.TLSSecrets = append(m.TLSSecrets, tls)

	// Validate manifest's TLS secret by creating a TLS secret.
	_, _, err := configmapandsecret.MakeTLSSecretAndGenerateName(tls)
	if err != nil {
		return err
	}

	return nil
}

func tlsSecretExists(m *manifest.Manifest, name string) bool {
	for _, s := range m.TLSSecrets {
		if name == s.Name {
			return true
		}
	}
	return false
}
