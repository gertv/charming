// Copyright 2015 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gertv/charming"
)

type TemplateConfig struct {
	Name       string `json:"name"`
	Stylesheet string `json:"stylesheet"`
	Source     string
}

func LoadTemplateConfigs(config charming.CharmingConfig) []TemplateConfig {
	glob := fmt.Sprintf("%s/**/template.json", config.TemplateDir)
	log.Printf("Loading templates from %s", glob)
	names, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}

	tcs := make([]TemplateConfig, len(names))
	for i, name := range names {
		log.Printf(" - adding template definition from %s", name)

		file, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}

		var tc TemplateConfig
		err = json.NewDecoder(file).Decode(&tc)
		if err != nil {
			log.Fatal(err)
		}
		tc.Source = name
		tcs[i] = tc
	}

	return tcs
}
