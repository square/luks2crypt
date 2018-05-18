// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package postimaging

import (
	"log"

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
func Run(opts *Opts) error {
	cryptServerInfo := &escrow.CryptServerInfo{
		Server: opts.Server,
		URI:    opts.URI,
	}
	cryptServerData := &escrow.CryptServerData{}

	// gather system data (hostname, username, serialnumber)
	sysinfo, sysinfoErr := hwinfo.Gather()
	if sysinfoErr != nil {
		return sysinfoErr
	}
	cryptServerData.Serialnum = sysinfo.SerialNum
	cryptServerData.Hostname = sysinfo.Hostname
	cryptServerData.Username = sysinfo.Username
	log.Printf("%+v\n", sysinfo)

	// create a new random password
	pass, errPass := password.New()
	if errPass != nil {
		return errPass
	}
	cryptServerData.Pass = pass
	log.Printf("generated new random password\n")

	// test if the current password works before performing destructive actions
	passWorks, passWorksErr := luks.PassWorks(opts.CurPass, opts.LuksDev)
	if passWorksErr != nil || passWorks == false {
		return passWorksErr
	}

	// start destructive operations
	// cache new password to disk
	cache, cacheErr := localcache.SaveToDisk(
		cryptServerData.Pass,
		"/etc/luks2crypt/crypt_recovery_key.json",
	)
	if cacheErr != nil {
		return cacheErr
	}
	log.Printf("luks password cached locally in: \"%s\"\n", cache)

	// escrew new password to crypt-server
	log.Printf(
		"escrowing system data to: %s\n",
		opts.Server+opts.URI,
	)
	postRes, postErr := cryptServerData.PostCryptServer(*cryptServerInfo)
	if postErr != nil {
		return postErr
	}
	if postRes.StatusCode != 200 {
		log.Fatalf("error posting to crypt-server: \"%+v\"\n", postRes)
	}
	log.Printf("escrowed data to crypt-server\n")

	// change luks admin password to new password
	setPassErr := luks.SetRecoveryPassword(opts.CurPass, cryptServerData.Pass,
		opts.LuksDev)
	if setPassErr != nil {
		return setPassErr
	}
	log.Printf("changed luks password")

	return nil
}
