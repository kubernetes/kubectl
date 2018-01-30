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

	"path"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/kinflate"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

const appname = "helloworld"

const manifestTemplate = `apiVersion: manifest.k8s.io/v1alpha1
kind: Manifest
metadata:
  name: helloworld
description: helloworld does useful stuff.
namePrefix: some-prefix
# Labels to add to all objects and selectors.
# These labels would also be used to form the selector for apply --prune
# Named differently than “labels” to avoid confusion with metadata for this object
objectLabels:
  app: helloworld
objectAnnotations:
  note: This is a example annotation
resources:
- deployment.yaml
- service.yaml
# There could also be configmaps in Base, which would make these overlays
configmaps: []
# There could be secrets in Base, if just using a fork/rebase workflow
secrets: []
recursive: true
`

const deploymentTemplate = `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: helloworld
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx
`

const serviceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: helloworld
  labels:
    app: helloworld
spec:
  ports:
    - port: 8888
  selector:
    app: helloworld
`

type initOptions struct {
}

// NewCmdInit make the init command.
func NewCmdInit(out, errOut io.Writer, fs fs.FileSystem) *cobra.Command {
	var o initOptions

	cmd := &cobra.Command{
		Use:   "init",
		Short: "TDB",
		Long:  "TBD",
		Example: `
		TBD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(cmd, args)
			if err != nil {
				return err
			}
			err = o.Complete(cmd, args)
			if err != nil {
				return err
			}
			return o.RunInit(cmd, out, errOut, fs)
		},
	}
	return cmd
}

// Validate validates init command.
func (o *initOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes init command.
func (o *initOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// RunKinflate runs init command (do real work).
func (o *initOptions) RunInit(cmd *cobra.Command, out, errOut io.Writer, fs fs.FileSystem) error {
	if _, err := fs.Stat("helloworld"); err == nil {
		return fmt.Errorf("%q already exists", appname)
	}
	err := fs.Mkdir(appname, 0744)
	if err != nil {
		return err
	}

	err = writefile(kinflate.KubeManifestFileName, []byte(manifestTemplate), fs)
	if err != nil {
		return err
	}
	err = writefile("deployment.yaml", []byte(deploymentTemplate), fs)
	if err != nil {
		return err
	}
	err = writefile("service.yaml", []byte(serviceTemplate), fs)
	return err
}

func writefile(name string, content []byte, fs fs.FileSystem) error {
	f, err := fs.Create(path.Join(appname, name))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}
