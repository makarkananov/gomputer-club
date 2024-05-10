package output

import (
	"fmt"
	"gomputerClub/internal/core/domain"
	"time"
)

// Console is a console output strategy.
type Console struct{}

// NewConsole creates a new Console.
func NewConsole() *Console {
	return &Console{}
}

// PrintTime prints the time in the format "15:04".
func (c *Console) PrintTime(t time.Time) {
	fmt.Println(t.Format("15:04"))
}

// PrintEvents prints the events depending on their type.
func (c *Console) PrintEvents(events []domain.Event) {
	for _, event := range events {
		id := event.GetID()
		fmt.Printf("%s %d ", event.GetTime().Format("15:04"), id)
		switch id {
		case 1, 3, 4, 11:
			clientEvent := event.(domain.ClientEvent)
			fmt.Println(clientEvent.ClientName)
		case 2, 12:
			tableEvent := event.(domain.TableEvent)
			fmt.Println(tableEvent.ClientName, tableEvent.TableNumber)
		case 13:
			errorEvent := event.(domain.ErrorEvent)
			fmt.Println(errorEvent.ErrorName)
		}
	}
}

// PrintRevenue prints the revenue and busy duration of the tables.
func (c *Console) PrintRevenue(tables []*domain.Table) {
	for i, table := range tables {
		fmt.Println(
			i+1,
			table.Revenue(),
			fmtDuration(table.BusyDuration()),
		)
	}
}

// fmtDuration formats the duration in the format "HH:MM".
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
