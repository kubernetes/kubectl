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

type addGenericSecret struct {
	// Name of secret (required)
	Name string
	// FileSources to derive the secret from (optional)
	FileSources []string
	// LiteralSources to derive the secret from (optional)
	LiteralSources []string
	// EnvFileSource to derive the secret from (optional)
	EnvFileSource string
}

// validate validates required fields are set to support structured generation.
func (a *addGenericSecret) Validate(args []string) error {
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

func newCmdAddSecretGeneric(errOut io.Writer) *cobra.Command {
	var config addGenericSecret
	cmd := &cobra.Command{
		Use:     "generic NAME [--type=string] [--from-file=[key=]source] [--from-literal=key1=value1]",
		Short:   "Adds a secret from a local file, directory or literal value.",
		Long:    "",
		Example: "",
		RunE: func(_ *cobra.Command, args []string) error {
			err := config.Validate(args)
			if err != nil {
				return err
			}

			if len(args) != 1 {
				return fmt.Errorf("error: exactly one NAME is required, got %d", len(args))
			}
			config.Name = args[0]

			// TODO(apelisse,droot): Do something with that config.

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&config.FileSources, "from-file", []string{}, "Key files can be specified using their file path, in which case a default name will be given to them, or optionally with a name and file path, in which case the given name will be used.  Specifying a directory will iterate each named file in the directory that is a valid secret key.")
	cmd.Flags().StringArrayVar(&config.LiteralSources, "from-literal", []string{}, "Specify a key and literal value to insert in secret (i.e. mykey=somevalue)")
	cmd.Flags().StringVar(&config.EnvFileSource, "from-env-file", "", "Specify the path to a file to read lines of key=val pairs to create a secret (i.e. a Docker .env file).")

	return cmd
}

type addTLSSecret struct {
	Name string
	Cert string
	Key  string
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
func newCmdAddSecretTLS(errOut io.Writer) *cobra.Command {
	var config addTLSSecret
	cmd := &cobra.Command{
		Use:     "tls NAME --cert=path/to/cert/file --key=path/to/key/file",
		Short:   "Adds a TLS secret.",
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

	cmd.Flags().StringVar(&config.Cert, "cert", "", "Path to PEM encoded public key certificate.")
	cmd.Flags().StringVar(&config.Key, "key", "", "Path to private key associated with given certificate.")

	return cmd
}

func NewCmdAddSecret(errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Adds a secret using specified subcommand",
	}
	cmd.AddCommand(newCmdAddSecretGeneric(errOut))
	cmd.AddCommand(newCmdAddSecretTLS(errOut))

	return cmd
}
