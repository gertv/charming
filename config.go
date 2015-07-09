// Copyright 2015 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package charming

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration for Charming server consists of a template directory,
// a work directory and a listening port for the server
type CharmingConfig struct {
	TemplateDir string `json:"templateDir"`
	WorkDir     string `json:"workDir"`
	Listen      string `json:"listen"`
}

// Print the config to the log
func (cc CharmingConfig) Log() {
	log.Printf("Using configuration:")
	log.Printf(" - template directory is %s", cc.TemplateDir)
	log.Printf(" - work directory is %s", cc.WorkDir)
	log.Printf(" - listen on %s", cc.Listen)
}

// Read config from a files.
func ReadConfig(config string) (result CharmingConfig, err error) {
	log.Printf("Reading engine configuration from %s", config)
	file, err := os.Open(config)
	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&result)
	return result, err
}
