package main

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	expectedArgs := []*regexp.Regexp{
		regexp.MustCompile("^--debug$"),
		regexp.MustCompile("^--advertise-client-urls$"),
		regexp.MustCompile("^our etcd url$"),
		regexp.MustCompile("^--listen-client-urls$"),
		regexp.MustCompile("^our etcd url$"),
		regexp.MustCompile("^--data-dir$"),
		regexp.MustCompile("^.+"),
	}
	numExpectedArgs := len(expectedArgs)
	numGivenArgs := len(os.Args) - 1

	if numGivenArgs < numExpectedArgs {
		fmt.Printf("Expected at least %d args, only got %d\n", numExpectedArgs, numGivenArgs)
		os.Exit(2)
	}

	for i, argRegexp := range expectedArgs {
		givenArg := os.Args[i+1]
		if !argRegexp.MatchString(givenArg) {
			fmt.Printf("Expected arg '%s' to match '%s'\n", givenArg, argRegexp.String())
			os.Exit(1)
		}
	}
	fmt.Println("Everything is dandy")
	fmt.Fprintln(os.Stderr, "serving insecure client requests on 127.0.0.1:2379")

	time.Sleep(10 * time.Minute)
}
