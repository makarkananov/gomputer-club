package service

import (
	"gomputerClub/internal/core/domain"
	"gomputerClub/internal/core/port"
	"sort"
)

// ComputerClubManager is a manager for handling events in a computer club
type ComputerClubManager struct {
	club   *domain.ComputerClub
	output port.OutputStrategy
}

// NewComputerClubManager creates a new ComputerClubManager with the given output strategy
func NewComputerClubManager(output port.OutputStrategy) *ComputerClubManager {
	return &ComputerClubManager{
		output: output,
	}
}

// SetClub sets the current club for the manager
func (ccm *ComputerClubManager) SetClub(club *domain.ComputerClub) {
	ccm.club = club
}

// HandleEvent handles the given event
func (ccm *ComputerClubManager) HandleEvent(event domain.Event) {
	ccm.club.Events = append(ccm.club.Events, event)

	switch event.GetID() {
	case 1: // Client tries to enter the club
		clientEvent := event.(domain.ClientEvent)
		ccm.handleClientComing(clientEvent)
	case 2: // Client tries to take a table
		tableEvent := event.(domain.TableEvent)
		ccm.handleTakingTable(tableEvent)
	case 3: // Client wants to wait for a table
		clientEvent := event.(domain.ClientEvent)
		ccm.handleWaiting(clientEvent)
	case 4: // Client leaves the club
		clientEvent := event.(domain.ClientEvent)
		ccm.handleLeaving(clientEvent)
	default:
		errEvent := domain.NewErrorEvent(event.GetTime(), "UnknownEvent")
		ccm.club.Events = append(ccm.club.Events, errEvent)
	}
}

// handleClientComing handles the event of a client trying to enter the club
func (ccm *ComputerClubManager) handleClientComing(event domain.ClientEvent) {
	// Client is already in the club
	if _, ok := ccm.club.Clients[event.ClientName]; ok {
		errEvent := domain.NewErrorEvent(event.GetTime(), "YouShallNotPass")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// club is closed
	if !(event.GetTime().After(ccm.club.OpenTime) && event.GetTime().Before(ccm.club.CloseTime)) {
		errEvent := domain.NewErrorEvent(event.GetTime(), "NotOpenYet")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// Client enters the club
	ccm.club.Clients[event.ClientName] = domain.NewClient(event.ClientName)
}

// handleTakingTable handles the event of a client trying to take a table
func (ccm *ComputerClubManager) handleTakingTable(event domain.TableEvent) {
	client, ok := ccm.club.Clients[event.ClientName]

	// Client is not in the club
	if !ok {
		errEvent := domain.NewErrorEvent(event.GetTime(), "ClientUnknown")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// Table is already taken
	if ccm.club.Tables[event.TableNumber-1].IsBusy() {
		errEvent := domain.NewErrorEvent(event.GetTime(), "PlaceIsBusy")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// Client already has a table, so we free it
	if client.Status == domain.Active {
		ccm.club.Tables[client.TableNumber-1].EndSession(event.GetTime(), ccm.club.CostPerHour)
	}

	// Client takes a table
	ccm.club.Tables[event.TableNumber-1].StartSession(event.GetTime())
	client.TableNumber = event.TableNumber
	client.Status = domain.Active
	ccm.club.Clients[event.ClientName] = client
}

// handleWaiting handles the event when a client wants to wait for a table
func (ccm *ComputerClubManager) handleWaiting(event domain.ClientEvent) {
	client, ok := ccm.club.Clients[event.ClientName]

	// Client is not in the club
	if !ok {
		errEvent := domain.NewErrorEvent(event.GetTime(), "ClientUnknown")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// Client already has a table
	if client.Status == domain.Active {
		errEvent := domain.NewErrorEvent(event.GetTime(), "YouHaveAPlace")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// Client is already waiting
	if client.Status == domain.Waiting {
		errEvent := domain.NewErrorEvent(event.GetTime(), "YouAreInTheQueue")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// There are available tables
	for _, t := range ccm.club.Tables {
		if !t.IsBusy() {
			errEvent := domain.NewErrorEvent(event.GetTime(), "ICanWaitNoLonger!")
			ccm.club.Events = append(ccm.club.Events, errEvent)
			return
		}
	}

	// Queue is full
	if len(ccm.club.Queue) == len(ccm.club.Tables) {
		newClientEvent := domain.NewClientEvent(11, event.GetTime(), domain.Outcome, event.ClientName)
		ccm.club.Events = append(ccm.club.Events, newClientEvent)
		return
	}

	// Client waits for a table
	client.Status = domain.Waiting
	ccm.club.Queue = append(ccm.club.Queue, client)
	ccm.club.Clients[event.ClientName] = client
}

// handleLeaving handles the event when a client leaves the club
func (ccm *ComputerClubManager) handleLeaving(event domain.ClientEvent) {
	client, ok := ccm.club.Clients[event.ClientName]

	// Client is not in the club
	if !ok {
		errEvent := domain.NewErrorEvent(event.GetTime(), "ClientUnknown")
		ccm.club.Events = append(ccm.club.Events, errEvent)
		return
	}

	// Client was waiting for a table
	if client.Status == domain.Waiting {
		for i, c := range ccm.club.Queue {
			if c.Name == client.Name {
				ccm.club.Queue = append(ccm.club.Queue[:i], ccm.club.Queue[i+1:]...)
				break
			}
		}
	}

	// If client had a table, it is now available
	if client.Status == domain.Active {
		ccm.club.Tables[client.TableNumber-1].EndSession(event.GetTime(), ccm.club.CostPerHour)

		if len(ccm.club.Queue) > 0 { // If there are clients waiting for a table, the first one gets it
			waitingClient := ccm.club.Queue[0]
			waitingClient.Status = domain.Active
			waitingClient.TableNumber = client.TableNumber
			ccm.club.Queue = ccm.club.Queue[1:]
			ccm.club.Tables[client.TableNumber-1].StartSession(event.GetTime())

			newTableEvent := domain.NewTableEvent(
				12,
				event.GetTime(),
				domain.Outcome,
				waitingClient.Name,
				waitingClient.TableNumber,
			)
			ccm.club.Events = append(ccm.club.Events, newTableEvent)
		}
	}

	delete(ccm.club.Clients, client.Name)
}

// EndDay ends the day in the club. It prints the results and resets the club
func (ccm *ComputerClubManager) EndDay() {
	// Sorting clients by name in alphabetical order
	keys := make([]string, 0, len(ccm.club.Clients))
	for key := range ccm.club.Clients {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Ending the session for each client
	for _, key := range keys {
		client := ccm.club.Clients[key]
		if client.Status == domain.Active {
			ccm.club.Tables[client.TableNumber-1].EndSession(ccm.club.CloseTime, ccm.club.CostPerHour)
		}

		clientEvent := domain.NewClientEvent(11, ccm.club.CloseTime, domain.Outcome, key)
		ccm.club.Events = append(ccm.club.Events, clientEvent)
	}

	// Printing the results
	ccm.output.PrintTime(ccm.club.OpenTime)
	ccm.output.PrintEvents(ccm.club.Events)
	ccm.output.PrintTime(ccm.club.CloseTime)
	ccm.output.PrintRevenue(ccm.club.Tables)

	// Resetting the club
	ccm.club.Clients = make(map[string]*domain.Client)
	ccm.club.Queue = make([]*domain.Client, 0)
	ccm.club.Events = make([]domain.Event, 0)
	ccm.club.Tables = make([]*domain.Table, len(ccm.club.Tables))
}
