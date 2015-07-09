// Copyright 2015 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gertv/charming"
	"github.com/gertv/charming/template"
)

type Request struct {
	uuid     string
	template template.TemplateConfig
	body     []byte
}

type TaskRequest struct {
	Uuid   string
	Input  string
	Style  string
	Output string
}

type Task struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
}

type Engine struct {
	Config charming.CharmingConfig

	database map[string]Task

	requests  chan Request
	responses chan Task
}

func (e Engine) AllTasks() []Task {
	js := make([]Task, len(e.database))
	i := 0
	for _, r := range e.database {
		js[i] = r
		i += 1
	}
	return js
}

func Setup(cc charming.CharmingConfig) *Engine {
	requests := make(chan Request)
	responses := make(chan Task)
	database := make(map[string]Task)

	engine := Engine{cc, database, requests, responses}
	go engine.handleResponses()
	go engine.handleRequests()
	return &engine
}

func (e *Engine) handleResponses() {
	log.Printf("Starting the reponse handler goroutine")
	for {
		select {
		case response, ok := <-e.responses:
			if ok {
				log.Printf("Request %s is now done", response.Uuid)
				e.database[response.Uuid] = response
			}
		}
	}
}

func (e *Engine) handleRequests() {
	log.Printf("Starting the request handler goroutine")
	for {
		select {
		case request, ok := <-e.requests:
			if ok {
				log.Printf("Let's handle %s", request.uuid)
				e.Handle(request)
			}
		}
	}
}

func (e Engine) Task(uuid string) Task {
	return e.database[uuid]
}

func (e Engine) TaskOutput(uuid string) []byte {
	pdf, err := ioutil.ReadFile(fmt.Sprintf("/home/gertv/Projects/WIP/charming/work/%s/output.pdf", uuid))
	if err != nil {
		log.Fatal(err)
	}

	return pdf
}

func (e Engine) Submit(tc template.TemplateConfig, body []byte) Task {
	uuid := strconv.FormatInt(time.Now().UnixNano(), 36)
	task := Task{uuid, "submitted"}
	e.database[uuid] = task

	e.requests <- Request{uuid, tc, body}

	return task
}

func (e Engine) Handle(r Request) {
	tempdir := filepath.Join(e.Config.WorkDir, r.uuid)
	err := os.Mkdir(tempdir, 0755)
	if err != nil {
		log.Printf("Unable to create work directory for request %s", r.uuid)
		e.responses <- Task{r.uuid, "failed"}
	}
	log.Printf("Executing request in %s", tempdir)

	input, err := os.Create(filepath.Join(tempdir, "input.html"))
	ioutil.WriteFile(input.Name(), r.body, 400)
	log.Printf("Creating temp file for input: %s", input.Name())

	output, err := os.Create(filepath.Join(tempdir, "output.pdf"))
	log.Printf("Will be writing to %s", output.Name())

	logfile, err := os.Create(filepath.Join(tempdir, "log.txt"))
	log.Printf("Prince log output will be writter to %s", logfile.Name())

	style := filepath.Join(filepath.Dir(r.template.Source), r.template.Stylesheet)
	log.Printf("Using stylesheet : %s", style)

	cmd := exec.Command("prince", "-s", style, "-i", "html", "-o", output.Name(), input.Name())

	cmd.Stdout = logfile
	cmd.Stderr = logfile
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	e.responses <- Task{r.uuid, "done"}
}
