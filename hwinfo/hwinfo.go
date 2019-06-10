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

type machine struct{}

type machineGetter interface {
	getSysSerialNumber() (string, error)
	getHostname() (string, error)
	getUsername() (string, error)
	isUser(string) bool
}

// SystemInfo stores sys serial number, hostname, and username
type SystemInfo struct {
	SerialNum, Hostname, Username string
}

// getSysSerialNumber returns the system serial number
func (i machine) getSysSerialNumber() (string, error) {
	dmi := dmidecode.New()

	err := dmi.Run()
	if err != nil {
		return "", err
	}

	byNameData, err := dmi.SearchByName("System Information")
	if err != nil {
		return "", err
	}

	return byNameData[0]["Serial Number"], nil
}

// getHostname returns the system hostname
func (i machine) getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// getUsername returns the current username
func (i machine) getUsername() (string, error) {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		return sudoUser, nil
	}

	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.Username, nil

}

// isUser returns true if we were run under user u
func (i machine) isUser(u string) bool {
	currentUser, _ := user.Current()

	if currentUser.Username == u {
		return true
	}
	return false
}

// GetInfo impliments the interface to gather system information such as
// system serial number, hostname, and current user
func getInfo(g machineGetter) (SystemInfo, error) {
	var err error
	i := SystemInfo{}

	i.Hostname, err = g.getHostname()
	if err != nil {
		return SystemInfo{}, err
	}

	i.Username, err = g.getUsername()
	if err != nil {
		return SystemInfo{}, err
	}

	if !g.isUser("root") {
		return SystemInfo{}, errors.New("not supported under current user please use sudo")
	}

	i.SerialNum, err = g.getSysSerialNumber()
	if err != nil {
		return SystemInfo{}, err
	}

	return i, nil
}

// Gather collects the system serial number, hostname, and current user
func Gather() (SystemInfo, error) {
	g := machine{}
	info, err := getInfo(g)

	if err != nil {
		return SystemInfo{}, err
	}

	return info, nil
}
