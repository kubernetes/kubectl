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

	"github.com/spf13/cobra"
)

type addConfigMap struct {
	// Name of configMap (required)
	Name string
	// FileSources to derive the configMap from (optional)
	FileSources []string
	// LiteralSources to derive the configMap from (optional)
	LiteralSources []string
	// EnvFileSource to derive the configMap from (optional)
	EnvFileSource string
}

// validate validates required fields are set to support structured generation.
func (a *addConfigMap) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("name must be specified once")
	}
	a.Name = args[0]
	if len(a.EnvFileSource) == 0 && len(a.FileSources) == 0 && len(a.LiteralSources) == 0 {
		return fmt.Errorf("at least from-env-file, or from-file or from-literal must be set")
	}
	if len(a.EnvFileSource) > 0 && (len(a.FileSources) > 0 || len(a.LiteralSources) > 0) {
		return fmt.Errorf("from-env-file cannot be combined with from-file or from-literal")
	}
	// TODO: Should we check if the path exists? if it's valid, if it's within the same (sub-)directory?
	return nil
}

func NewCmdAddConfigMap(errOut io.Writer) *cobra.Command {
	var config addConfigMap
	cmd := &cobra.Command{
		Use:     "configmap NAME [--from-file=[key=]source] [--from-literal=key1=value1]",
		Short:   "Adds a configmap to your manifest file.",
		Long:    "",
		Example: "",
		RunE: func(_ *cobra.Command, args []string) error {
			err := config.Validate(args)
			if err != nil {
				return err
			}

			// TODO(apelisse,droot): Do something with that config.

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&config.FileSources, "from-file", []string{}, "Key file can be specified using its file path, in which case file basename will be used as configmap key, or optionally with a key and file path, in which case the given key will be used.  Specifying a directory will iterate each named file in the directory whose basename is a valid configmap key.")
	cmd.Flags().StringArrayVar(&config.LiteralSources, "from-literal", []string{}, "Specify a key and literal value to insert in configmap (i.e. mykey=somevalue)")
	cmd.Flags().StringVar(&config.EnvFileSource, "from-env-file", "", "Specify the path to a file to read lines of key=val pairs to create a configmap (i.e. a Docker .env file).")

	return cmd
}
