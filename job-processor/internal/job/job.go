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
	Execute   func() error // actual work logic
}

// New creates a new job with timestamp
func New(id string, priority Priority, execute func() error) *Job {
	return &Job{
		ID:        id,
		Priority:  priority,
		CreatedAt: time.Now(),
		Execute:   execute,
	}
}

// String makes debugging easier.
func (j *Job) String() string {
	return fmt.Sprintf("[Job ID=%s Priority=%s]", j.ID, j.Priority)
}
