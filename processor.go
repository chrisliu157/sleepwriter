package main

import (
	"encoding/json"
	"net/http"

	render "github.com/chrisliu156/sleepwriter/http"
	"github.com/gorilla/mux"
)

func ProcessGetJob(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobId := vars["jobId"]
	if jobId == "" {
		render.RenderError(rw, http.StatusBadRequest, render.Error.Validation, "No job ID found in URI path.")
		return
	}

	job, getErr := sys.store.Get(jobId)
	if getErr != nil {
		render.RenderError(rw, http.StatusNotFound, render.Error.NotFound, "Could not find job.")
		return
	}

	rw.Write(job)
}

func ProcessSleepWriteRequest(rw http.ResponseWriter, req *http.Request) {
	body, err := render.ParseBody(rw, req.Body)
	if err != nil {
		render.RenderError(rw, http.StatusBadRequest, render.Error.Validation, "Could not parse request.")
		return
	}

	var writeReq SleepWriter
	if err = json.Unmarshal(body, &writeReq); err != nil {
		render.RenderError(rw, http.StatusBadRequest, render.Error.Validation, "Could not parse request.")
		return
	}

	valid := writeReq.IsValid()
	if valid != nil {
		render.RenderError(rw, http.StatusBadRequest, render.Error.Validation, valid.Error())
		return
	}

	jobId, asyncErr := writeReq.Async()
	if asyncErr != nil {
		render.RenderError(rw, http.StatusInternalServerError, render.Error.InternalServer, "Please try again later.")
		return
	}

	sys.log.Debugf("Submitted Job Id: %v", jobId)
	render.RenderAsync(rw, jobId)
	return
}
