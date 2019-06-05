package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	args := os.Args[0:1]
	args = append(args, "-version")

	err := run(args)

	if err != nil {
		t.Errorf("error when calling run: %v", err)
	}
}
