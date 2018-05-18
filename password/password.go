// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package password

import (
	"errors"
	"log"
	"strings"

	"github.com/sethvargo/go-diceware/diceware"
)

// password holds the password(s) for this package
type password struct {
	new string
}

// properties contains info used when generating new passwords
type properties struct {
	numWords      int
	wordSeperator string
}

// generate creates a unique and sufficently long set of chars to be
// used as a password
func generate(prop *properties) string {
	list, err := diceware.Generate(prop.numWords)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Join(list, prop.wordSeperator)
}

// validate checks if password is valid
func (p password) validate(prop *properties) error {
	list := strings.Split(p.new, prop.wordSeperator)
	if len(list) != prop.numWords {
		return errors.New("generated password does not match requested length")
	}
	return nil
}

// New creates a new password and validates it
func New() (string, error) {
	props := &properties{numWords: 5, wordSeperator: "."}
	pass := &password{new: generate(props)}
	validationErr := pass.validate(props)
	if validationErr != nil {
		return "", validationErr
	}
	return pass.new, nil
}
