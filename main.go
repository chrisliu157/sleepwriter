package main

import (
	"net/http"
	"os"
	"time"

	"github.com/chrisliu156/sleepwriter/store"
	"github.com/chrisliu156/sleepwriter/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	logLevel     = utils.GetEnv("log_level", "DEBUG")
	jobQueueName = utils.GetEnv("job_queue_name", "job_queue")
)

type System struct {
	store store.Store
	log   *logrus.Logger
}

var sys *System

func init() {
	sys = &System{}
	sys.log = logrus.New()
	sys.log.Out = os.Stdout

	switch logLevel {
	case "DEBUG":
		sys.log.Level = logrus.DebugLevel
	case "INFO":
		sys.log.Level = logrus.InfoLevel
	}

	pool, poolErr := store.NewStore()
	if poolErr != nil {
		sys.log.Fatalf("Initialization - Error - Store connection error - %v", poolErr)
	}
	sys.store = *pool
}

func setRoutes() *mux.Router {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/jobs/{jobId}", func(w http.ResponseWriter, r *http.Request) {
		ProcessGetJob(w, r)
	}).Methods("GET")

	v1.HandleFunc("/sleep_writers", func(w http.ResponseWriter, r *http.Request) {
		ProcessSleepWriteRequest(w, r)
	}).Methods("POST")

	return router
}

func main() {

	sys.log.Infoln("Starting Sleep Writer Service..")
	sys.log.Info("Ready to serve HTTP requests at 0.0.0.0:3000...")

	server := &http.Server{
		Handler:      setRoutes(),
		Addr:         "0.0.0.0:3000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	sys.log.Fatal(server.ListenAndServe())
}
