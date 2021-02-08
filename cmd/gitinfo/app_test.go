package main_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	cmd "github.com/pgmig/gitinfo/cmd/gitinfo"
)

func TestRunErrors(t *testing.T) {
	// Save original args
	a := os.Args

	tests := []struct {
		name string
		code int
		args []string
	}{
		{"Help", 3, []string{"-h"}},
		{"UnknownFlag", 2, []string{"-0"}},
		{"Path is empty", 1, []string{}},
		{"Cannot write file", 1, []string{"--gi.file", ".", "."}},
		{"File is not a dir (skip)", 0, []string{"main.go"}},
	}
	for _, tt := range tests {
		os.Args = append([]string{a[0]}, tt.args...)
		var c int
		cmd.Run(func(code int) { c = code })
		assert.Equal(t, tt.code, c, tt.name)
	}
	// Restore original args
	os.Args = a
}

func TestRun(t *testing.T) {
	// Save original args
	a := os.Args
	os.Args = append([]string{a[0]}, "--debug",
		"./",
	)
	var c int
	cmd.Run(func(code int) { c = code })
	assert.Equal(t, 0, c, "Normal run")
	// Restore original args
	os.Args = a
}
