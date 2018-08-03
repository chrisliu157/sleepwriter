package main

import (
	"fmt"
	"time"

	"github.com/chrisliu156/sleepwriter/jobs"
)

type SleepWriter struct {
	SleepDuration int    `json:"sleep_duration"`
	Value         string `json:"value"`
	Priority      int    `json:"priority"`
}

func (writer SleepWriter) IsValid() error {
	if writer.SleepDuration == 0 {
		return fmt.Errorf("Sleep duration cannot be empty")
	}

	if writer.Value == "" {
		return fmt.Errorf("Value cannot be empty")
	}

	return nil
}

func (writer SleepWriter) Async() (string, error) {
	job, _ := jobs.NewSleepWriterJob(time.Duration(writer.SleepDuration)*time.Second, writer.Value, writer.Priority, jobQueueName)
	jobErr := job.Delay(sys.store)
	if jobErr != nil {
		return "", jobErr
	}
	return job.JobId, nil
}
