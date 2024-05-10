package output

import (
	"gomputerClub/internal/core/domain"
	"time"
)

type MockOutputStrategy struct {
	PrintTimeCalls    []time.Time
	PrintEventsCalls  [][]domain.Event
	PrintRevenueCalls [][]*domain.Table
}

func (m *MockOutputStrategy) PrintTime(t time.Time) {
	m.PrintTimeCalls = append(m.PrintTimeCalls, t)
}

func (m *MockOutputStrategy) PrintEvents(events []domain.Event) {
	m.PrintEventsCalls = append(m.PrintEventsCalls, events)
}

func (m *MockOutputStrategy) PrintRevenue(tables []*domain.Table) {
	m.PrintRevenueCalls = append(m.PrintRevenueCalls, tables)
}
