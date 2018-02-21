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
	"errors"
	"io"

	"github.com/spf13/cobra"

	"k8s.io/kubectl/pkg/kinflate/tree"
	"k8s.io/kubectl/pkg/kinflate/util"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
	"k8s.io/utils/exec"
)

type diffOptions struct {
	manifestPath string
}

// newCmdDiff makes the diff command.
func newCmdDiff(out, errOut io.Writer, fs fs.FileSystem) *cobra.Command {
	var o diffOptions

	cmd := &cobra.Command{
		Use:     "diff",
		Short:   "diff between transformed resources and untransformed resources",
		Long:    "diff between transformed resources and untransformed resources and the subpackages are all transformed.",
		Example: `diff -f .`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(cmd, args)
			if err != nil {
				return err
			}
			err = o.Complete(cmd, args)
			if err != nil {
				return err
			}
			return o.RunDiff(out, errOut, fs)
		},
	}

	cmd.Flags().StringVarP(&o.manifestPath, "filename", "f", "", "Pass in a Kube-manifest.yaml file or a directory that contains the file.")
	cmd.MarkFlagRequired("filename")
	return cmd
}

// Validate validates diff command.
func (o *diffOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.New("The diff command takes no arguments.")
	}
	return nil
}

// Complete completes diff command.
func (o *diffOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// RunInit writes a manifest file.
func (o *diffOptions) RunDiff(out, errOut io.Writer, fs fs.FileSystem) error {
	printer := util.Printer{}
	diff := util.DiffProgram{
		Exec:   exec.New(),
		Stdout: out,
		Stderr: errOut,
	}

	inflateOp := inflateOptions{manifestPath: o.manifestPath, mode: tree.ModeNormal}
	kobj1, err := inflateOp.runInflate(fs)
	if err != nil {
		return err
	}
	transformedDir, err := util.WriteToDir(kobj1, "transformed", printer)
	if err != nil {
		return err
	}
	defer transformedDir.Delete()

	inflateNoOp := inflateOptions{manifestPath: o.manifestPath, mode: tree.ModeNoop}
	kobj2, err := inflateNoOp.runInflate(fs)
	if err != nil {
		return err
	}
	noopDir, err := util.WriteToDir(kobj2, "noop", printer)
	if err != nil {
		return err
	}
	defer noopDir.Delete()

	return diff.Run(noopDir.Name, transformedDir.Name)
}
