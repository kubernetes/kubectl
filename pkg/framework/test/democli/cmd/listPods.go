// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// listPodsCmd represents the listPods command
var listPodsCmd = &cobra.Command{
	Use:   "listPods",
	Short: "List all pods",
	Long:  `Give a list of all pods known by the system`,
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, err := cmd.Flags().GetString("api-url")
		if err != nil {
			panic(err)
		}
		runGetPods(apiURL)
	},
}

func runGetPods(apiURL string) {
	config, err := clientcmd.BuildConfigFromFlags(apiURL, "")
	if err != nil {
		panic(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	if len(pods.Items) > 0 {
	} else {
		fmt.Println("There are no pods.")
	}
}

func init() {
	RootCmd.AddCommand(listPodsCmd)

	// Here you will define your flags and configuration settings.

	listPodsCmd.Flags().String("api-url", "http://localhost:8080", "URL of the APIServer to connect to")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listPodsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listPodsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
