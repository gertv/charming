// Copyright 2015 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gertv/charming"
	"github.com/gertv/charming/api"
	"github.com/gertv/charming/engine"
	"github.com/gertv/charming/template"
)

// Main function for the charmingd executable.
//
// Really, no summary?
func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: charming <template directory>\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Template directory has not been set")
		os.Exit(1)
	}

	config, err := charming.ReadConfig(args[0])
	if err == nil {
		config.Log()
	} else {
		log.Fatalf("Unable to parse configuration file: %s", err)
	}

	e := engine.Setup(config)
	tcs := template.LoadTemplateConfigs(config)

	api := api.Setup(e, tcs)
	http.Handle("/", api.Router)

	log.Fatal(http.ListenAndServe(":6060", nil))
}
