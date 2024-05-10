package domain

import "time"

// ComputerClub represents a computer club with clients, tables, events, and a queue.
type ComputerClub struct {
	Clients     map[string]*Client
	Tables      []*Table
	Events      []Event
	Queue       []*Client
	OpenTime    time.Time
	CloseTime   time.Time
	CostPerHour int
}

// NewComputerClub creates a new computer club.
func NewComputerClub(openTime, closeTime time.Time, numberOfTables int, costPerHour int) *ComputerClub {
	cc := &ComputerClub{
		OpenTime:    openTime,
		CloseTime:   closeTime,
		CostPerHour: costPerHour,
		Clients:     make(map[string]*Client),
	}

	cc.Tables = make([]*Table, 0, numberOfTables)
	for range numberOfTables {
		cc.Tables = append(cc.Tables, NewTable())
	}

	return cc
}
