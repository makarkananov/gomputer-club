package domain

import "time"

// EventType represents the type of event that occurred.
type EventType uint8

const (
	Income EventType = iota
	Outcome
)

// Event represents an event in the computer club.
type Event interface {
	GetID() int
	GetTime() time.Time
	GetType() EventType
}

// BaseEvent represents the base fields of an event.
type BaseEvent struct {
	ID        int
	Time      time.Time
	EventType EventType
}

func (be BaseEvent) GetID() int {
	return be.ID
}

func (be BaseEvent) GetTime() time.Time {
	return be.Time
}

func (be BaseEvent) GetType() EventType {
	return be.EventType
}

// ClientEvent represents an event related to a client.
type ClientEvent struct {
	BaseEvent
	ClientName string
}

// NewClientEvent creates a new client event.
func NewClientEvent(id int, time time.Time, eventType EventType, clientName string) ClientEvent {
	return ClientEvent{
		BaseEvent: BaseEvent{
			ID:        id,
			Time:      time,
			EventType: eventType,
		},
		ClientName: clientName,
	}
}

// ErrorEvent represents an error event.
type ErrorEvent struct {
	BaseEvent
	ErrorName string
}

// NewErrorEvent creates a new error event.
func NewErrorEvent(time time.Time, errorName string) ErrorEvent {
	return ErrorEvent{
		BaseEvent: BaseEvent{
			ID:        13,
			Time:      time,
			EventType: Outcome,
		},
		ErrorName: errorName,
	}
}

// TableEvent represents an event related to a table.
type TableEvent struct {
	ClientEvent
	TableNumber int
}

// NewTableEvent creates a new table event.
func NewTableEvent(id int, time time.Time, eventType EventType, clientName string, tableNumber int) TableEvent {
	return TableEvent{
		ClientEvent: ClientEvent{
			BaseEvent: BaseEvent{
				ID:        id,
				Time:      time,
				EventType: eventType,
			},
			ClientName: clientName,
		},
		TableNumber: tableNumber,
	}
}
