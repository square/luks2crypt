// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package localcache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type expectedData struct {
	newPass, cacheFile string
}

// createTempFile allocates a temporary file in the system $TMPDIR
func createTempFile(t *testing.T, dir string, file string) *os.File {
	tmpfile, err := ioutil.TempFile(dir, file)
	if err != nil {
		t.Fatalf("error creating temp file %v", err)
	}

	t.Logf("created tmp file '%v'", tmpfile.Name())

	return tmpfile
}

// verifyTestFile validates the expected vs actual results written to the test
// file
func verifyTestFile(file string, expected expectedData, t *testing.T) error {
	actual := &CacheData{}

	file = path.Clean(file)

	data, err := os.Open(file)
	defer data.Close()
	if err != nil {
		t.Errorf("error reading file contents %v", err)
	}
	byteData, err := ioutil.ReadAll(data)
	if err != nil {
		t.Errorf("error converting file to byte array %v", err)
	}

	json.Unmarshal(byteData, &actual)

	if actual.AdminPassNew != expected.newPass {
		t.Errorf("expected '%v' password written to disk but got '%v'", expected.newPass, actual.AdminPassNew)
	}
	return nil
}

func TestSaveToDisk(t *testing.T) {
	file := createTempFile(t, "", "go-TestSaveToDisk-*.yaml")
	defer os.Remove(file.Name())

	expected := expectedData{
		newPass:   "1234.foo.bar.test",
		cacheFile: file.Name(),
	}

	actual, err := SaveToDisk(expected.newPass, expected.cacheFile)
	if err != nil {
		t.Errorf("error occurred writing file: %v", err)
	}

	err = verifyTestFile(actual, expected, t)
	if err != nil {
		t.Errorf("error verifying file contents %v", err)
	}
}
