package domain

import (
	"math"
	"time"
)

// Table represents a table in the computer club.
type Table struct {
	isBusy       bool
	busyDuration time.Duration
	revenue      int
	lastStart    time.Time
}

// NewTable creates a new table.
func NewTable() *Table {
	return &Table{}
}

// StartSession starts a session at the table.
func (t *Table) StartSession(time time.Time) {
	t.isBusy = true
	t.lastStart = time
}

// EndSession ends a session at the table. It calculates the revenue and the busy duration.
func (t *Table) EndSession(time time.Time, costPerHour int) {
	if !t.isBusy {
		return
	}

	t.isBusy = false
	sessionDuration := time.Sub(t.lastStart)
	t.busyDuration += sessionDuration
	t.revenue += int(math.Ceil(sessionDuration.Hours())) * costPerHour
}

func (t *Table) Revenue() int {
	return t.revenue
}

func (t *Table) IsBusy() bool {
	return t.isBusy
}

func (t *Table) BusyDuration() time.Duration {
	return t.busyDuration
}
