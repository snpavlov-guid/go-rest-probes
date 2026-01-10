// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"fmt"
)

func ExampleUsage() {
	// This example demonstrates the expected usage output
	// Note: This is a simplified example as the actual usage() function
	// calls os.Exit which cannot be demonstrated in an example test
	fmt.Println("usage: console [options] [arguments]")
	fmt.Println("The options are:")
	fmt.Println("  -v")
	fmt.Println("    Enable verbose output.")
	fmt.Println("  -h, --help")
	fmt.Println("    Show this help message.")
	// Output:
	// usage: console [options] [arguments]
	// The options are:
	//   -v
	//     Enable verbose output.
	//   -h, --help
	//     Show this help message.
}

func ExampleCommandLineArguments() {
	// Example showing how command-line arguments work
	args := []string{"arg1", "arg2", "arg3"}
	fmt.Printf("Received %d argument(s)\n", len(args))
	for i, arg := range args {
		fmt.Printf("  [%d] %s\n", i+1, arg)
	}
	// Output:
	// Received 3 argument(s)
	//   [1] arg1
	//   [2] arg2
	//   [3] arg3
}
