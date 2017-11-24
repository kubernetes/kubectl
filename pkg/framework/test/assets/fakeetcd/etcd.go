package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	expectedArgs := []string{
		"--debug",
		"--advertise-client-urls",
		"our etcd url",
		"--listen-client-urls",
		"our etcd url",
		"--data-dir",
	}
	numExpectedArgs := len(expectedArgs)
	numGivenArgs := len(os.Args) - 1

	if numGivenArgs < numExpectedArgs {
		fmt.Printf("Expected at least %d args, only got %d\n", numExpectedArgs, numGivenArgs)
		os.Exit(2)
	}

	for i, arg := range expectedArgs {
		givenArg := os.Args[i+1]
		if arg != givenArg {
			fmt.Printf("Expected arg %s, got arg %s\n", arg, givenArg)
			os.Exit(1)
		}
	}
	fmt.Println("Everything is dandy")

	time.Sleep(10 * time.Minute)
}
