// Copyright 2015 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gertv/charming/engine"
	"github.com/gertv/charming/template"
	"github.com/gorilla/mux"
)

type Api struct {
	engine *engine.Engine
	Router *mux.Router
}

type TemplateInfo struct {
	Name      string `json:"name"`
	SubmitUrl string `json:"submitUrl"`
}

func newTemplateInfo(tc template.TemplateConfig) TemplateInfo {
	return TemplateInfo{tc.Name, fmt.Sprintf("http://localhost:3000/submit/%s", tc.Name)}
}

type TaskInfo struct {
	engine.Task
	OutputUrl string `json:"outputUrl"`
}

func Setup(e *engine.Engine, tcs []template.TemplateConfig) Api {
	r := mux.NewRouter()

	r.HandleFunc("/task", getAllJobs(e)).Methods("GET")
	r.HandleFunc("/task/{uuid}", getJob(e)).Methods("GET")
	r.HandleFunc("/task/{uuid}/output.pdf", getOutput(e)).Methods("GET")

	r.HandleFunc("/template", getAllTemplates(tcs)).Methods("GET")
	for _, tc := range tcs {
		uri := fmt.Sprintf("/submit/%s", tc.Name)
		log.Printf("Starting template handler at %s", uri)
		r.HandleFunc(uri, executeTemplate(e, tc))
	}

	return Api{e, r}
}

func getAllTemplates(tcs []template.TemplateConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tis := make([]TemplateInfo, len(tcs))
		for i, tc := range tcs {
			tis[i] = newTemplateInfo(tc)
		}

		json, err := json.Marshal(tis)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(json)
	}
}

func executeTemplate(e *engine.Engine, tc template.TemplateConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		result := e.Submit(tc, body)

		http.Redirect(w, r, fmt.Sprintf("/task/%s", result.Uuid), 303)
		defer r.Body.Close()
	}
}

func getAllJobs(e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := e.AllTasks()

		json, err := json.Marshal(wrapTasks(jobs))
		if err != nil {
			log.Fatal(err)
		}

		w.Write(json)
	}
}

func wrapTasks(input []engine.Task) []TaskInfo {
	result := make([]TaskInfo, len(input))

	for i, t := range input {
		result[i] = newTaskInfo(t)
	}

	return result
}

func newTaskInfo(input engine.Task) TaskInfo {
	return TaskInfo{input, fmt.Sprintf("http://localhost:6060/task/%s/output.pdf", input.Uuid)}
}

func getJob(e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]
		job := e.Task(uuid)

		json, err := json.Marshal(newTaskInfo(job))
		if err != nil {
			log.Fatal(err)
		}

		w.Write(json)
	}
}

func getOutput(e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]
		data := e.TaskOutput(uuid)

		w.Header().Set("Content-Type", "application/pdf")
		w.Write(data)
	}
}
