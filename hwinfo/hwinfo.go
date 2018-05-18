// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package hwinfo

import (
	"errors"
	"os"
	"os/user"

	"github.com/dselans/dmidecode"
)

// SystemInfo stores sys serial number, hostname, and username
type SystemInfo struct {
	SerialNum, Hostname, Username string
}

// getSysSerialNumber returns the system serial number
func getSysSerialNumber() (string, error) {
	dmi := dmidecode.New()

	dmiRunErr := dmi.Run()
	if dmiRunErr != nil {
		return "", dmiRunErr
	}

	byNameData, byNameErr := dmi.SearchByName("System Information")
	if byNameErr != nil {
		return "", byNameErr
	}

	return byNameData["Serial Number"], nil
}

// getHostname returns the system hostname
func getHostname() (string, error) {
	return os.Hostname()
}

// getUsername returns the current username
func getUsername() (string, error) {
	username, err := user.Current()
	return username.Username, err
}

// isRootUser returns true if we were run under root
func (sysinfo SystemInfo) isRootUser() bool {
	if sysinfo.Username == "root" {
		return true
	}
	return false
}

// Gather collects the system serial number, hostname, and current user
func Gather() (*SystemInfo, error) {
	sysinfo := &SystemInfo{}

	hostname, hostnameErr := getHostname()
	if hostnameErr != nil {
		return nil, hostnameErr
	}
	sysinfo.Hostname = hostname

	username, usernameErr := getUsername()
	if usernameErr != nil {
		return nil, usernameErr
	}
	sysinfo.Username = username

	if !sysinfo.isRootUser() {
		return nil, errors.New("not supported under current user please use sudo")
	}

	serial, serialErr := getSysSerialNumber()
	if serialErr != nil {
		return nil, serialErr
	}
	sysinfo.SerialNum = serial

	return sysinfo, nil
}
