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

	"github.com/spf13/cobra"
)

type initOptions struct {
}

// NewCmdInit make the init command.
func NewCmdInit(out, errOut io.Writer) *cobra.Command {
	var o initOptions

	cmd := &cobra.Command{
		Use:   "init",
		Short: "TDB",
		Long:  "TBD",
		Example: `
		TBD`,
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
			err = o.RunInit(cmd, out, errOut)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				os.Exit(1)
			}
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
func (o *initOptions) RunInit(cmd *cobra.Command, out, errOut io.Writer) error {
	_, err := out.Write([]byte("Hello I am init.\n"))
	return err
}
