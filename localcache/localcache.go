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
func (data CacheData) mkConfDir() error {
	data.Path = path.Clean(data.Path)
	confPath := path.Dir(data.Path)

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		log.Printf("creating config and cache dir: %s\n", confPath)
		err = os.Mkdir(confPath, os.FileMode(data.PathMode))
		if err != nil {
			return err
		}
	}

	return nil
}

// marshalJSONData converts json data into a marshal byte blob for writting
func (data CacheData) marshalJSONData() ([]byte, error) {
	marshal, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

// SaveToDisk saves the luks recovery password on disk
// returns errors if any
func SaveToDisk(newPass string, cacheFile string) (string, error) {
	data := CacheData{
		AdminPassNew: newPass,
		Path:         cacheFile,
		PathMode:     0700,
	}

	marshalData, err := data.marshalJSONData()
	if err != nil {
		return "", err
	}

	err = data.mkConfDir()
	if err != nil {
		return "", err
	}

	cacheFile = path.Clean(cacheFile)
	err = ioutil.WriteFile(cacheFile, marshalData, 0600)
	if err != nil {
		return "", err
	}
	return cacheFile, nil
}
