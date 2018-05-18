// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package escrow

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

// CryptServerInfo is used to create an object with info about the escrow server
type CryptServerInfo struct {
	Server, URI string
}

// CryptServerData stores the data to be escrowed
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

	encodeErr := encoder.Encode(data, form)
	if encodeErr != nil {
		return nil, encodeErr
	}

	client := new(http.Client)
	res, postErr := client.PostForm(cryptServer, form)
	if postErr != nil {
		return nil, postErr
	}
	return res, nil
}
