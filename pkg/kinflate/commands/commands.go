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
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

// NewDefaultCommand returns the default (aka root) command for kinflate command.
func NewDefaultCommand() *cobra.Command {
	fsys := fs.MakeRealFS()
	stdOut, stdErr := os.Stdout, os.Stderr

	c := &cobra.Command{
		Use:   "kinflate",
		Short: "kinflate manages Kubernetes application configuration.",
		Long: `
kinflate manages Kubernetes application configuration.

Find more information at:
	https://github.com/kubernetes/kubectl/tree/master/cmd/kinflate
`,
	}

	f := newFactory()
	f.bindFlags(c.PersistentFlags())

	c.AddCommand(
		newCmdInflate(f, stdOut, stdErr, fsys),
		newCmdDiff(stdOut, stdErr, fsys),
		newCmdInit(stdOut, stdErr, fsys),
		// 'add' sub command
		newCmdAdd(stdOut, stdErr, fsys),
		// 'set' sub command
		newCmdSet(stdOut, stdErr, fsys),
	)

	return c
}

// newAddCommand returns an instance of 'add' subcommand.
func newCmdAdd(stdOut, stdErr io.Writer, fsys fs.FileSystem) *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Adds configmap/resource/secret to the manifest.",
		Long:  "",
		Example: `
	# Adds a configmap to the manifest
	kinflate add configmap NAME --from-literal=k=v

	# Adds a secret to the manifest
	kinflate add secret NAME --from-literal=k=v

	# Adds a resource to the manifest
	kinflate add resource <filepath>
`,
		Args: cobra.MinimumNArgs(1),
	}
	c.AddCommand(
		newCmdAddResource(stdOut, stdErr, fsys),
		newCmdAddConfigMap(stdErr, fsys),
	)
	return c
}

// newSetCommand returns an instance of 'set' subcommand.
func newCmdSet(stdOut, stdErr io.Writer, fsys fs.FileSystem) *cobra.Command {
	c := &cobra.Command{
		Use:   "set",
		Short: "Sets the value of different fields in manifest.",
		Long:  "",
		Example: `
	# Sets the nameprefix field
	kinflate set nameprefix <prefix-value>
`,
		Args: cobra.MinimumNArgs(1),
	}

	c.AddCommand(
		newCmdSetNamePrefix(stdOut, stdErr, fsys),
	)
	return c
}
