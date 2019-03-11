// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
)

// This is a simple webserver that listens on 127.0.0.1:443 and prints form data
// Uses the snake oil certs and is useful in debuging luks2crypt escrow form data
func main() {
	http.HandleFunc("/checkin/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		for key, value := range r.Form {
			fmt.Printf("%s = %s\n", key, value)
		}
	})

	err := http.ListenAndServeTLS(
		":443",
		"/etc/ssl/certs/ssl-cert-snakeoil.pem",
		"/etc/ssl/private/ssl-cert-snakeoil.key",
		nil,
	)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
