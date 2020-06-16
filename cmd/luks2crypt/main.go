// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/square/luks2crypt/pkg/postimaging"

	"golang.org/x/crypto/ssh/terminal"
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
				cli.StringFlag{Name: "authuser, u",
					Usage: "Basic auth username for Crypt server."},
				cli.StringFlag{Name: "authpass, P",
					Usage: "Basic auth password for Crypt server. If omitted and authuser/u is specified, the password will be prompted for on the terminal."},
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

// optPostImaging sets up, rotates, and escrows the initial encryption password
func optPostImaging(c *cli.Context) error {
	cryptURL := "https://" + c.String("cryptserver")
	opts := postimaging.Opts{
		LuksDev:  c.String("luksdevice"),
		CurPass:  c.String("currentpassword"),
		Server:   cryptURL,
		URI:      c.String("cryptendpoint"),
		AuthUser: c.String("authuser"),
		AuthPass: c.String("authpass"),
	}
	if (opts.AuthUser != "") && (opts.AuthPass == "") {
		fmt.Printf("Password: ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			err = fmt.Errorf("error getting basic auth password: %v", err)
			return cli.NewExitError(err, 1)
		}
		opts.AuthPass = password
	}
	err := postimaging.Run(opts)
	if err != nil {
		err = fmt.Errorf("error setting escrow passphrase: %v", err)
		return cli.NewExitError(err, 1)
	}
	return nil
}

// main calls run to execute luks2crypt
func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
