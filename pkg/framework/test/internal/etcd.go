package internal

import "fmt"

func MakeEtcdArgs(input DefaultedProcessInput) []string {
	args := []string{
		"--debug",
		"--listen-peer-urls=http://localhost:0",
		fmt.Sprintf("--advertise-client-urls=%s", input.URL.String()),
		fmt.Sprintf("--listen-client-urls=%s", input.URL.String()),
		fmt.Sprintf("--data-dir=%s", input.Dir),
	}
	return args
}
