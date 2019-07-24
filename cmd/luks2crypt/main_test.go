// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"testing"
)

var testCases = []struct {
	arg  string
	want error
}{
	{
		arg:  "-version",
		want: nil,
	},
}

func TestMain(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("flag %s should produce error of %s", tc.arg, tc.want), func(t *testing.T) {
			args := os.Args[0:1]
			args = append(args, tc.arg)

			err := run(args)

			if err != tc.want {
				t.Errorf("error when calling run: %v", err)
			}
		})
	}
}
