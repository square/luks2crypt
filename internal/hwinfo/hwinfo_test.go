// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package hwinfo

import (
	"testing"
)

type FakeSystemInfo struct {
	SerialNum, Hostname, Username string
	isRoot                        bool
}

func (i FakeSystemInfo) getSysSerialNumber() (string, error) {
	return i.SerialNum, nil
}

func (i FakeSystemInfo) getHostname() (string, error) {
	return i.Hostname, nil
}

func (i FakeSystemInfo) getUsername() (string, error) {
	return i.Username, nil
}

func (i FakeSystemInfo) isUser(u string) bool {
	return i.isRoot
}

func TestGetInfo(t *testing.T) {
	expected := FakeSystemInfo{
		SerialNum: "1234FooBar",
		Hostname:  "testing.example.com",
		Username:  "testinguser",
		isRoot:    true,
	}

	actual, err := getInfo(expected)

	if err != nil {
		t.Errorf("failed to get system information: %v", err)
	}

	if actual.SerialNum != expected.SerialNum {
		t.Errorf("expected serialnumber %v, instead got %v", expected.SerialNum, actual.SerialNum)
	}

	if actual.Hostname != expected.Hostname {
		t.Errorf("expected hostname %v, instead got %v", expected.Hostname, actual.Hostname)
	}

	if actual.Username != expected.Username {
		t.Errorf("expected username %v, instead got %v", expected.Username, actual.Username)
	}
}

func TestNonrootGetInfo(t *testing.T) {
	var expectedEmpty SystemInfo
	expected := FakeSystemInfo{
		Username: "testinguser",
		isRoot:   false,
	}

	actual, err := getInfo(expected)

	if err == nil {
		t.Errorf("checking for a non-root user should have returned an error")
	}

	if actual != expectedEmpty {
		t.Errorf("expected actual to be empty. Got %v", actual)
	}
}
