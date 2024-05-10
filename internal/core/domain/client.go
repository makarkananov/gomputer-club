package domain

type ClientStatus int

const (
	Inactive ClientStatus = iota // Client is doing nothing
	Waiting                      // Client is waiting for a table
	Active                       // Client has a table
)

// Client represents a client of the computer club.
type Client struct {
	Name        string
	Status      ClientStatus
	TableNumber int // -1 if client doesn't have a table
}

// NewClient creates a new client.
func NewClient(name string) *Client {
	return &Client{
		Name:        name,
		Status:      Inactive,
		TableNumber: -1,
	}
}
