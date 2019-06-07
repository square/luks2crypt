// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
)

type cryptServerHandler struct{}

// This is a simple webserver that listens on 127.0.0.1:443 and prints form data
// Uses the snake oil certs and is useful in debuging luks2crypt escrow form data
func (c *cryptServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	for key, value := range r.Form {
		log.Printf("%s = %s\n", key, value)
	}
}

func main() {
	mux := http.NewServeMux()
	serv := &cryptServerHandler{}

	mux.Handle("/checkin/", serv)

	log.Println("Listening...")
	err := http.ListenAndServeTLS(
		":8443",
		"/etc/ssl/certs/ssl-cert-snakeoil.pem",
		"/etc/ssl/private/ssl-cert-snakeoil.key",
		mux,
	)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
