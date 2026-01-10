// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Console is a simple console application demonstrating
// how to write a basic command-line program.
//
// Usage:
//
//	console [options] [arguments]
//
// The options are:
//
//	-v
//		Enable verbose output.
//
//	-h, --help
//		Show this help message.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	verbose = flag.Bool("v", false, "Enable verbose output")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: console [options] [arguments]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	// Configure logging for a command-line program.
	log.SetFlags(0)
	log.SetPrefix("console: ")

	// Parse flags.
	flag.Usage = usage
	flag.Parse()

	// Get remaining arguments.
	args := flag.Args()

	// Display help if requested.
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		usage()
		return
	}

	// Run actual logic.
	if *verbose {
		log.Println("Verbose mode enabled")
		fmt.Printf("Received %d argument(s):\n", len(args))
		for i, arg := range args {
			fmt.Printf("  [%d] %s\n", i+1, arg)
		}
		return
	}

	if len(args) > 0 {
		fmt.Println("Arguments:", args)
	} else {
		fmt.Println("Console application started. Use -v for verbose output or -h for help.")
	}
}
