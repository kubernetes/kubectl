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
	"io"

	"github.com/ghodss/yaml"

	"errors"

	"github.com/spf13/cobra"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/constants"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

type setPrefixNameOptions struct {
	prefix string
}

// NewCmdSetPrefixName sets the value of the namePrefix field in the manifest.
func NewCmdSetPrefixName(out, errOut io.Writer, fsys fs.FileSystem) *cobra.Command {
	var o setPrefixNameOptions

	cmd := &cobra.Command{
		Use:   "setprefixname",
		Short: "Sets the value of the namePrefix field in the manifest.",
		Long:  "Sets the value of the namePrefix field in the manifest.",
		//
		Example: `
The command
  setprefixname acme-
will add the field "namePrefix: acme-" to the manifest file if it doesn't exist,
and overwrite the value with "acme-" if the field does exist.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(args)
			if err != nil {
				return err
			}
			err = o.Complete(cmd, args)
			if err != nil {
				return err
			}
			return o.RunSetPrefixName(out, errOut, fsys)
		},
	}
	return cmd
}

// Validate validates setPrefixName command.
func (o *setPrefixNameOptions) Validate(args []string) error {
	if len(args) != 1 {
		return errors.New("must specify exactly one prefix value")
	}
	// TODO: add further validation on the value.
	o.prefix = args[0]
	return nil
}

// Complete completes setPrefixName command.
func (o *setPrefixNameOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// RunSetPrefixName runs setPrefixName command (does real work).
func (o *setPrefixNameOptions) RunSetPrefixName(out, errOut io.Writer, fsys fs.FileSystem) error {
	content, err := fsys.ReadFile(constants.KubeManifestFileName)
	if err != nil {
		return err
	}

	// TODO: Refactor manifest reading to a common location.
	// See pkg/kinflate/util.go:loadManifestPkg
	var m manifest.Manifest
	err = yaml.Unmarshal(content, &m)
	if err != nil {
		return err
	}

	m.NamePrefix = o.prefix

	bytes, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	err = fsys.WriteFile(constants.KubeManifestFileName, bytes)
	if err != nil {
		return err
	}
	return nil
}
