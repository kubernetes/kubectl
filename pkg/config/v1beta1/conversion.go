/*
Copyright 2025 The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"

	conversion "k8s.io/apimachinery/pkg/conversion"
	config "k8s.io/kubectl/pkg/config"
)

func Convert_config_AllowlistEntry_To_v1beta1_AllowlistEntry(in *config.AllowlistEntry, out *AllowlistEntry, s conversion.Scope) error {
	in.Command = out.Command

	return nil
}

// The internal AllowlistEntry type does not have the `Name` field, which is deprecated as of v1.36. Convert `Name` to `Command` where possible, and return an error if both `Name` and `Command` are supplied.
func Convert_v1beta1_AllowlistEntry_To_config_AllowlistEntry(in *AllowlistEntry, out *config.AllowlistEntry, s conversion.Scope) error {
	if len(in.Name) != 0 && len(in.Command) != 0 {
		return fmt.Errorf("both `Name` and `Command` were supplied. `Name` is deprecated, use `Command` instead")
	}

	// when both `Name` and `Command` are empty, propagate the empty value and
	// allow validation to catch it later since it's a validation error, not a
	// conversion error
	cmd := in.Command
	if len(in.Name) != 0 {
		cmd = in.Name
	}

	out.Command = cmd
	return nil
}
