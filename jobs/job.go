package jobs

import "github.com/chrisliu156/sleepwriter/store"

const (
	PENDING    = "Pending"
	INPROGRESS = "In Progress"
	SUCCESS    = "Success"
)

type Job interface {
	Delay(store.Store) error
	Process(db store.Store) error
	GetPriority() int
	GetStatus() string
	SetStatus(status string)
	GetId() string
}
