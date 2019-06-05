// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package password

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	expected := properties{
		numWords:      5,
		wordSeperator: ".",
	}

	actual, err := New()

	if err != nil {
		t.Errorf("error generating password: %v", err)
	}

	actualWords := strings.Split(actual, expected.wordSeperator)

	if len(actualWords) != expected.numWords {
		t.Errorf("expected %v words but only generated %v", expected.numWords, len(actualWords))
	}
}
