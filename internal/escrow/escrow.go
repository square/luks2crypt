// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package escrow

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

// CryptServerInfo is used to create an object with info about the escrow server
type CryptServerInfo struct {
	Server, URI, Username, Password string
}

// CryptServerData stores the data to be escrowed
// cryptserver expects the following form data recovery_password, serial,
// macname, username
// see: https://github.com/grahamgilbert/Crypt-Server/blob/master/server/views.py#L442
type CryptServerData struct {
	Pass      string `schema:"recovery_password"`
	Serialnum string `schema:"serial"`
	Hostname  string `schema:"macname"`
	Username  string `schema:"username"`
}

// PostCryptServer submits the luks password and machine info to crypt-server
func (data CryptServerData) PostCryptServer(escrowServer CryptServerInfo) (*http.Response, error) {
	cryptServer := escrowServer.Server + escrowServer.URI

	encoder := schema.NewEncoder()
	form := url.Values{}

	err := encoder.Encode(data, form)
	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	req, err := http.NewRequest("POST", cryptServer, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if (escrowServer.Username != "") && (escrowServer.Password != "") {
		req.SetBasicAuth(escrowServer.Username, escrowServer.Password)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
