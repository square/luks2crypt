// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	serv := &cryptServerHandler{}
	expected := url.Values{}
	expected.Add("test", "data")

	req, err := http.NewRequest(
		"GET",
		"/checkin/",
		strings.NewReader(expected.Encode()),
	)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serv.ServeHTTP)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned %v expected %v", status, http.StatusOK)
	}
}
