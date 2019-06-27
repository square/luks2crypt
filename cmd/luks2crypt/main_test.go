// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

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
