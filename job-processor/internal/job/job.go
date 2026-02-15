package job

import (
	"fmt"
	"time"
)

// Job representas
type Job struct {
	ID        string       // unique identifier
	Priority  Priority     // HIGH or NORMAL
	CreatedAt time.Time    // tracking
	Run       func() error // actual work logic
}

// New creates a new job with timestamp
func New(id string, priority Priority, run func() error) *Job {
	return &Job{
		ID:        id,
		Priority:  priority,
		CreatedAt: time.Now(),
		Run:       run,
	}
}

// Execute safely runs the job logic.
func (j *Job) Execute() error {
	if j.Run == nil {
		return nil
	}
	return j.Run()
}

// String makes debugging easier.
func (j *Job) String() string {
	return fmt.Sprintf("[Job ID=%s Priority=%s]", j.ID, j.Priority)
}
