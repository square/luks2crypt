// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package luks

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/partition/mbr"
)

type testDisk struct {
	path, pass, newPass string
	size                int64
	luksVersion         int
}

// createTempDir allocates a temporary directory in the system $TMPDIR
func createTempDir(t *testing.T, dir string, prefix string) string {
	tmpdir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		t.Fatalf("error creating temp dir %v", err)
	}

	t.Logf("created tmp file '%v'", tmpdir)

	return tmpdir
}

// create allocates and formats a disk to run tests against
func (d testDisk) create(t *testing.T) {
	luksDev := &Settings{
		NewPass:     d.pass,
		LuksDevice:  d.path,
		LuksVersion: d.luksVersion,
	}

	disk, err := diskfs.Create(d.path, d.size, diskfs.SectorSizeDefault)
	if err != nil {
		t.Fatalf("error creating test disk: %v", err)
	}

	table := &mbr.Table{
		LogicalSectorSize:  512,
		PhysicalSectorSize: 512,
	}

	err = disk.Partition(table)
	if err != nil {
		t.Errorf("error partitioning test filesystem %v", err)
	}

	_, err = formatSetPassword(luksDev.NewPass, luksDev.LuksDevice, luksDev.LuksVersion)
	if err != nil {
		t.Errorf("error creating test luks device: %v", err)
	}
}

func TestPassWorks(t *testing.T) {
	dir := createTempDir(t, "", "go-TestPassWorks")
	defer os.RemoveAll(dir)

	expected := testDisk{
		path:        path.Clean(dir + "/luksdisk.img"),
		size:        int64(10 * 1024 * 1024), // 10MB
		pass:        "testPassw0rd!",
		luksVersion: 1,
	}
	expected.create(t)

	_, err := PassWorks(expected.pass, expected.path, expected.luksVersion)
	if err != nil {
		t.Errorf("error checking if '%v' is the password for '%v'. Got %v",
			expected.pass,
			expected.path,
			err,
		)
	}
}

func TestSetRecoveryPassword(t *testing.T) {
	dir := createTempDir(t, "", "go-TestSetRecoveryPassword")
	defer os.RemoveAll(dir)

	expected := testDisk{
		path:        path.Clean(dir + "/luksdisk.img"),
		size:        int64(10 * 1024 * 1024), // 10MB
		pass:        "testPassw0rd!",
		newPass:     "Th!sIsTh3NewPassw0d*",
		luksVersion: 1,
	}
	expected.create(t)

	err := SetRecoveryPassword(expected.pass, expected.newPass, expected.path, expected.luksVersion)
	if err != nil {
		t.Errorf("error changing password from '%v' to '%v' on '%v'. Got %v",
			expected.pass,
			expected.newPass,
			expected.path,
			err,
		)
	}
}
