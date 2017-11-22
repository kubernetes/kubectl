package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	expectedArgs := []string{
		"--advertise-client-urls",
		"our etcd url",
		"--data-dir",
		"our data directory",
		"--listen-client-urls",
		"our etcd url",
		"--debug",
	}

	for i, arg := range os.Args[1:] {
		if arg != expectedArgs[i] {
			fmt.Printf("Expected arg %s, got arg %s", expectedArgs[i], arg)
			os.Exit(1)
		}
	}
	fmt.Println("Everything is dandy")

	time.Sleep(10 * time.Minute)
}
