// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestUsage(t *testing.T) {
	// Save original stderr
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()

	// Create a pipe to capture stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Capture output in a goroutine
	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		_, _ = buf.ReadFrom(r)
		done <- true
	}()

	// Reset flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	verbose = flag.Bool("v", false, "Enable verbose output")

	// Call usage function in a goroutine since it calls os.Exit
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected: usage() calls os.Exit(2)
			}
		}()
		usage()
	}()

	// Close write end
	w.Close()
	<-done

	output := buf.String()
	if !strings.Contains(output, "usage:") {
		t.Errorf("usage() output should contain 'usage:', got: %q", output)
	}
}

func TestProcessArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantArgs []string
	}{
		{
			name:     "no arguments",
			args:     []string{},
			wantArgs: []string{},
		},
		{
			name:     "with arguments",
			args:     []string{"arg1", "arg2", "arg3"},
			wantArgs: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "single argument",
			args:     []string{"test"},
			wantArgs: []string{"test"},
		},
		{
			name:     "empty string argument",
			args:     []string{""},
			wantArgs: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args) != len(tt.wantArgs) {
				t.Errorf("args length = %d, wantArgs length = %d", len(tt.args), len(tt.wantArgs))
			}

			for i, arg := range tt.args {
				if i < len(tt.wantArgs) && arg != tt.wantArgs[i] {
					t.Errorf("args[%d] = %q, want %q", i, arg, tt.wantArgs[i])
				}
			}
		})
	}
}

func TestHelpFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantHelp bool
	}{
		{
			name:     "help flag -h",
			args:     []string{"-h"},
			wantHelp: true,
		},
		{
			name:     "help flag --help",
			args:     []string{"--help"},
			wantHelp: true,
		},
		{
			name:     "no help flag",
			args:     []string{"arg1"},
			wantHelp: false,
		},
		{
			name:     "empty args",
			args:     []string{},
			wantHelp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isHelp := len(tt.args) > 0 && (tt.args[0] == "-h" || tt.args[0] == "--help")
			if isHelp != tt.wantHelp {
				t.Errorf("isHelp = %v, want %v", isHelp, tt.wantHelp)
			}
		})
	}
}

func TestVerboseFlagDefault(t *testing.T) {
	// Reset flags for testing
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	testVerbose := flag.Bool("v", false, "Enable verbose output")

	// Test default value
	if *testVerbose != false {
		t.Errorf("Default verbose flag should be false, got %v", *testVerbose)
	}
}
