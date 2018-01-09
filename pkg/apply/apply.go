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

package apply

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
)

type ApplyOptions struct {
	ConfigFile string
}

// NewCmdApply creates a new apply command.
func NewCmdApply(in io.Reader, out, errOut io.Writer) *cobra.Command {
	var o ApplyOptions

	cmd := &cobra.Command{
		Use:   "apply -f [path]",
		Short: "Declaratively create or update a resource using a config file",
		Long:  "Declaratively create of update a resource using a config file",
		Example: `
		# Run apply on a config file".
		apply -f FILENAME`,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Validate(cmd, args)
			if err != nil {
				panic(err)
			}
			err = o.RunApply(cmd, in, out, errOut)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&o.ConfigFile, "filename", "f", "", "Pass in config.yaml file.")
	return cmd
}

// Validate validates apply command. "args" should be empty.
func (o *ApplyOptions) Validate(cmd *cobra.Command, args []string) error {

	// Validate parameters
	if cmd == nil {
		return errors.New("Missing cobra command for apply command")
	}
	if len(args) != 0 {
		return errors.New("Unexpected args")
	}

	return nil
}

// RunApply runs apply command (do real work).
func (o *ApplyOptions) RunApply(cmd *cobra.Command, in io.Reader, out, errOut io.Writer) error {
	return nil
}
