package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"local/notorious/logging"
	"local/notorious/opts"
	"os"
	"strings"
)

func main() {
	o, err := opts.Parse()
	log := logging.Error(os.Stderr)
	if err != nil {
		log.Fatalf("parsing command-line options: %v. Try notorious --help for more information on command-line flags.", err)
	}
	err = Notorious(o, os.Stdin, os.Stdout, os.Stderr)

	if err != nil {
		log.Fatalf("running notorious: %v", err)
	}
	os.Exit(0)
}

// when main() is called, stdin, stdout, and stderr are what you'd expect them to be
func Notorious(o opts.Opts, stdin io.Reader, stdout, stderr io.Writer) error {
	// this will choke on large files and might waste a lot of memory.
	// ideally, we'd use a buffer of some kind and only read in enough to handle 
	// our lines of context before and after
	b, err := ioutil.ReadAll(stdin)
	if err != nil {
		return err
	}
	logger := logging.Debug(stderr, o.Verbose)
	lines := strings.Split(string(b), "\n")

	toWrite := make([]bool, len(lines))

	for i, line := range lines {
		if o.Matches(line) {
			logger.Printf("match: line %d: %s", i, line)
			start := findMax(0, i-o.Context.Before)
			end := findMin(len(lines), i+o.Context.After+1)
			for i := start; i < end; i++ {
				toWrite[i] = true

			}
		}
	}
	
	// go doesn't buffer stdin, out, or err by default.
	// we're about to do a bunch of repeated writes very quickly, so it's quicker to buffer them.
	out := bufio.NewWriter(stdout)
	if o.LineNumbers {
		for i, line := range lines {
			if toWrite[i] {
				fmt.Fprintf(out, "%d\t%s\n", i, line)
			}
		}
	} else {
		for i, line := range lines {
			if toWrite[i] {
				fmt.Fprintln(out, line)
			}
		}
	}

	// call Flush() when we're done, or some of the writes might never make it to 
	// the underlying io.Writer (in most cases, the OS)
	return out.Flush()
}

func findMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func findMin(a, b int) int {
	if a > b {
		return b
	}
	return a
}
