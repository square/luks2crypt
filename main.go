// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/square/luks2crypt/postimaging"

	cli "gopkg.in/urfave/cli.v1"
)

var (
	// VERSION set during build
	// go build -ldflags "-X main.VERSION=1.2.3"
	VERSION = "0.0.1"
)

// run setups up cli arg handling and executes luks2crypt
func run(args []string) error {
	app := cli.NewApp()
	app.Name = "luks2crypt"
	app.Usage = "Generates a random luks password, escrows it, and rotates slot 0 on root."
	app.Version = VERSION

	app.Commands = []cli.Command{
		{
			Name:   "version",
			Usage:  "show the application version",
			Action: optVersion,
		},
		{
			Name:   "postimaging",
			Usage:  "Changes the public luks password post imaging to a unique passphrase",
			Action: optPostImaging,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "luksdevice, l",
					Usage: "Luks Device to rotate password on. Ex. /dev/sda3"},
				cli.StringFlag{Name: "currentpassword, p",
					Usage: "Password to unlock and update device"},
				cli.StringFlag{Name: "cryptserver, s",
					Usage: "Crypt Server to escrow recovery key to. Ex. cryptserver.example.com"},
				cli.StringFlag{Name: "cryptendpoint, e",
					Usage: "The Crypt Server endpoint to use when escrowing keys",
					Value: "/checkin/"},
			},
		},
	}

	return app.Run(args)
}

// optVersion returns the application version. Typically, this is the git sha
func optVersion(c *cli.Context) error {
	fmt.Println(VERSION)
	return nil
}

// optPostImaging sets up, rotates, and escrows the inital encryption password
func optPostImaging(c *cli.Context) error {
	cryptURL := "https://" + c.String("cryptserver")
	opts := &postimaging.Opts{
		LuksDev: c.String("luksdevice"),
		CurPass: c.String("currentpassword"),
		Server:  cryptURL,
		URI:     c.String("cryptendpoint"),
	}
	err := postimaging.Run(opts)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

// main calls run to execute luks2crypt
func main() {
	err := run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
