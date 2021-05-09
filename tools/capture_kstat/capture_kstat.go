package main

// Serializes all Named kstats under the given path, to disk. Useful for generating fixture data
// to mock out tests which require kstats.
// The data is of type []*kstat.Named

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/siebenmann/go-kstat"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: capture-kstat <kstat> <file>")
		os.Exit(1)
	}

	kstatName := os.Args[1]
	file := os.Args[2]
	chunks := strings.Split(kstatName, ":")

	if len(chunks) != 3 {
		fmt.Fprintln(os.Stderr, "kstat must be of the form module:instance:name")
		os.Exit(1)
	}

	module := chunks[0]
	instance, _ := strconv.Atoi(chunks[1])
	name := chunks[2]

	token, err := kstat.Open()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get kstat token.")
	}

	rawKstat, err := token.Lookup(module, instance, name)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get kstat.")
		os.Exit(2)
	}

	stats, err := rawKstat.AllNamed()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get named kstat data.")
		os.Exit(2)
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err = enc.Encode(stats)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not encode data %v\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(file, buf.Bytes(), 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write serialized data to disk: %v\n", err)
		os.Exit(1)
	}
}
