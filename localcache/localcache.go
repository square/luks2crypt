// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package localcache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// CacheData stores the luks cache data to disk
type CacheData struct {
	AdminPassNew string `json:"admin_pass_new"`
	Path         string `json:"-"`
	PathMode     int    `json:"-"`
}

// mkConfDir creates the configuration dir. Typically, /etc/luks2crypt
func (data *CacheData) mkConfDir() error {
	confPath := path.Dir(data.Path)

	if _, existErr := os.Stat(confPath); os.IsNotExist(existErr) {
		log.Printf("creating config and cache dir: %s\n", confPath)
		mkErr := os.Mkdir(confPath, os.FileMode(data.PathMode))
		if mkErr != nil {
			return mkErr
		}
	}

	return nil
}

// marshalJSONData converts json data into a marshal byte blob for writting
func (data *CacheData) marshalJSONData() ([]byte, error) {
	marshal, marshalErr := json.MarshalIndent(data, "", "\t")
	if marshalErr != nil {
		return nil, marshalErr
	}
	return marshal, nil
}

// SaveToDisk saves the luks recovery password on disk
// returns errors if any
func SaveToDisk(newPass string, cacheFile string) (string, error) {
	data := &CacheData{
		AdminPassNew: newPass,
		Path:         cacheFile,
		PathMode:     0700,
	}

	marshalData, marshalDataErr := data.marshalJSONData()
	if marshalDataErr != nil {
		return "", marshalDataErr
	}

	mkdirErr := data.mkConfDir()
	if mkdirErr != nil {
		return "", mkdirErr
	}

	writeErr := ioutil.WriteFile(cacheFile, marshalData, 0600)
	if writeErr != nil {
		return "", writeErr
	}
	return cacheFile, nil
}
