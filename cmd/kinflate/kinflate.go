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

package main

import (
	"os"

	"github.com/spf13/cobra"

	"k8s.io/kubectl/pkg/kinflate/commands"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

func main() {
	var cmd = &cobra.Command{}
	cmd.AddCommand(commands.NewCmdInflate(os.Stdout, os.Stderr))
	cmd.AddCommand(commands.NewCmdInit(os.Stdout, os.Stderr, fs.MakeRealFS()))
	cmd.AddCommand(commands.NewCmdAddConfigMap(os.Stderr))

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
