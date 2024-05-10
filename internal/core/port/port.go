package port

import (
	"gomputerClub/internal/core/domain"
	"time"
)

// OutputStrategy is a strategy for outputting club data
type OutputStrategy interface {
	PrintTime(t time.Time)
	PrintEvents(events []domain.Event)
	PrintRevenue(tables []*domain.Table)
}

// ComputerClubManager represents a manager for a computer club
type ComputerClubManager interface {
	SetClub(club *domain.ComputerClub)
	HandleEvent(event domain.Event)
	EndDay()
}
