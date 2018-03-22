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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/kinflate/app"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
	"k8s.io/kubectl/pkg/loader"
)

type inflateOptions struct {
	manifestPath string
}

// newCmdInflate creates a new inflate command.
func newCmdInflate(f *factory, out, errOut io.Writer, fs fs.FileSystem) *cobra.Command {
	var o inflateOptions

	cmd := &cobra.Command{
		Use:   "inflate -f [path]",
		Short: "Use a Manifest file to generate a set of api resources",
		Long:  "Use a Manifest file to generate a set of api resources",
		Example: `
		# Use the Kube-manifest.yaml file under somedir/ to generate a set of api resources.
		inflate -f somedir/`,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Validate(cmd, args)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				os.Exit(1)
			}
			err = o.Complete(cmd, args)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				os.Exit(1)
			}
			err = o.RunInflate(f, out, errOut, fs)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&o.manifestPath, "filename", "f", "", "Pass in a Kube-manifest.yaml file or a directory that contains the file.")
	cmd.MarkFlagRequired("filename")
	return cmd
}

// Validate validates inflate command.
func (o *inflateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes inflate command.
func (o *inflateOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// RunInflate runs inflate command (do real work).
func (o *inflateOptions) RunInflate(f *factory, out, errOut io.Writer, fs fs.FileSystem) error {
	_, err := f.openAPISchema()
	if err != nil {
		return err
	}

	l := loader.Init([]loader.SchemeLoader{loader.NewFileLoader(fs)})

	absPath, err := filepath.Abs(o.manifestPath)
	if err != nil {
		return err
	}

	rootLoader, err := l.New(absPath)
	if err != nil {
		return err
	}

	application, err := app.New(rootLoader)
	if err != nil {
		return err
	}

	allResources, err := application.Resources()

	if err != nil {
		return err
	}

	// Output the objects.
	res, err := kutil.Encode(allResources)
	if err != nil {
		return err
	}
	_, err = out.Write(res)
	return err
}
