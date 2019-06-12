// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package postimaging

import (
	"fmt"
	"log"
	"net/http"

	"github.com/square/luks2crypt/escrow"
	"github.com/square/luks2crypt/hwinfo"
	"github.com/square/luks2crypt/localcache"
	"github.com/square/luks2crypt/luks"
	"github.com/square/luks2crypt/password"
)

// Opts is used to store the options needed for postimaging functions
type Opts struct {
	LuksDev, CurPass, Server, URI string
}

// Run post imaging password creation, set, and escrow. Returns an error
func Run(opts Opts) error {
	cryptServerInfo := escrow.CryptServerInfo{
		Server: opts.Server,
		URI:    opts.URI,
	}
	cryptServerData := escrow.CryptServerData{}

	// gather system data (hostname, username, serialnumber)
	sysinfo, err := hwinfo.Gather()
	if err != nil {
		return err
	}
	cryptServerData.Serialnum = sysinfo.SerialNum
	cryptServerData.Hostname = sysinfo.Hostname
	cryptServerData.Username = sysinfo.Username
	log.Printf("%+v\n", sysinfo)

	// create a new random password
	pass, err := password.New()
	if err != nil {
		return err
	}
	cryptServerData.Pass = pass
	log.Println("generated new random password")

	// test if the current password works before performing destructive actions
	passWorks, err := luks.PassWorks(opts.CurPass, opts.LuksDev)
	if err != nil || passWorks == false {
		return err
	}

	// start destructive operations
	// cache new password to disk
	cache, err := localcache.SaveToDisk(
		cryptServerData.Pass,
		"/etc/luks2crypt/crypt_recovery_key.json",
	)
	if err != nil {
		return err
	}
	log.Printf("luks password cached locally in: '%s'\n", cache)

	// escrew new password to crypt-server
	log.Printf(
		"escrowing system data to: %s\n",
		opts.Server+opts.URI,
	)
	postRes, err := cryptServerData.PostCryptServer(cryptServerInfo)
	if err != nil {
		return err
	}
	if postRes.StatusCode != http.StatusOK {
		return fmt.Errorf("error posting to crypt-server: '%+v'", postRes)
	}
	log.Println("escrowed data to crypt-server")

	// change luks admin password to new password
	err = luks.SetRecoveryPassword(opts.CurPass, cryptServerData.Pass,
		opts.LuksDev)
	if err != nil {
		return err
	}
	log.Println("changed luks passphrase")

	return nil
}
