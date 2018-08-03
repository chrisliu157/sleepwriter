package jobs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chrisliu156/sleepwriter/store"
	uuid "github.com/satori/go.uuid"
)

type SleepWriterJob struct {
	SleepDuration time.Duration `json:"sleep_duration"`
	Value         string        `json:"value"`
	CreatedAt     time.Time     `json:"created_at"`
	Priority      int           `json:"priority"`
	Status        string        `json:"status"`
	JobId         string        `json:"job_id"`
	Queue         string        `json:"queue"`
}

func NewSleepWriterJob(duration time.Duration, value string, priority int, queue string) (*SleepWriterJob, error) {
	sleepWriter := new(SleepWriterJob)
	sleepWriter.Priority = priority
	sleepWriter.Value = value
	sleepWriter.SleepDuration = duration
	sleepWriter.Queue = queue
	return sleepWriter, nil
}

func (writer *SleepWriterJob) Delay(db store.Store) error {
	writer.JobId = fmt.Sprintf("job:%s", uuid.Must(uuid.NewV4()).String())
	writer.CreatedAt = time.Now()
	writer.Status = PENDING
	payload, _ := json.Marshal(writer)

	err := db.PQueueAdd(writer.Queue, writer.Priority, writer.JobId, payload)
	return err
}

func (writer SleepWriterJob) Process(db store.Store) error {
	time.Sleep(writer.SleepDuration)
	key := strings.Replace(writer.JobId, "job:", "value:", -1)
	insertErr := db.Set(key, []byte(writer.Value))
	return insertErr
}

func (writer SleepWriterJob) GetPriority() int {
	return writer.Priority
}

func (writer SleepWriterJob) GetStatus() string {
	return writer.Status
}

func (writer *SleepWriterJob) SetStatus(status string) {
	writer.Status = status
}

func (writer *SleepWriterJob) SetCreatedAt(createdAt time.Time) {
	writer.CreatedAt = createdAt
}

func (writer SleepWriterJob) GetId() string {
	return writer.JobId
}
