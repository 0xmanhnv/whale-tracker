package listeners

import "database/sql"

// Scheduler data structure
type Scheduler struct {
	db        *sql.DB
	listeners Listeners
}

// NewScheduler creates a new scheduler
func NewScheduler(db *sql.DB, listeners Listeners) Scheduler {
	return Scheduler{
		db:        db,
		listeners: listeners,
	}
}
