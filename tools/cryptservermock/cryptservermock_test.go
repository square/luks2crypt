// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var testCases = []struct {
	data url.Values
	uri  string
	want int
}{
	{
		uri: "/checkin/",
		data: url.Values{
			"test": []string{"data"},
		},
		want: http.StatusOK,
	},
}

func TestMain(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s, expect: %v", tc.data, tc.want), func(t *testing.T) {
			serv := &cryptServerHandler{}

			req, err := http.NewRequest(
				"GET",
				tc.uri,
				strings.NewReader(tc.data.Encode()),
			)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(serv.ServeHTTP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.want {
				t.Errorf("handler returned %v expected %v", status, tc.want)
			}
		})
	}
}
