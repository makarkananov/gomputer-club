package service

import (
	"gomputerClub/internal/adapters/output"
	"gomputerClub/internal/core/domain"
	"testing"
	"time"
)

func TestComputerClubManager_HandleClientComing_ClubClosed(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	event := domain.NewClientEvent(1, startTime.Add(-1*time.Minute), domain.Income, "John")
	manager.handleClientComing(event)
	if _, exists := manager.club.Clients["John"]; exists {
		t.Errorf("Expected client %q to not be added when the club is closed", "John")
	}

	event = domain.NewClientEvent(1, startTime.Add(9*time.Hour), domain.Income, "John")
	manager.handleClientComing(event)
	if _, exists := manager.club.Clients["John"]; exists {
		t.Errorf("Expected client %q to not be added when the club is closed", "John")
	}
}

func TestComputerClubManager_HandleTakingTable_TableTaken(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	comingEvent := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.handleClientComing(comingEvent)
	comingEvent = domain.NewClientEvent(1, startTime.Add(2*time.Minute), domain.Income, "Mike")
	manager.handleClientComing(comingEvent)

	tableEvent := domain.NewTableEvent(2, startTime.Add(3*time.Minute), domain.Income, "John", 1)
	manager.handleTakingTable(tableEvent)

	tableEvent2 := domain.NewTableEvent(2, startTime.Add(4*time.Minute), domain.Income, "Mike", 1)
	manager.handleTakingTable(tableEvent2)

	if manager.club.Clients["John"].TableNumber != 1 {
		t.Errorf("Expected John to have table number 1")
	}
	if manager.club.Clients["Mike"].Status != domain.Inactive {
		t.Errorf("Expected Mike's status to be inactive")
	}
	if !manager.club.Tables[0].IsBusy() {
		t.Errorf("Expected table 1 to be busy")
	}
}

func TestComputerClubManager_HandleClientComing(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	event := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.HandleEvent(event)

	if _, exists := manager.club.Clients["John"]; !exists {
		t.Errorf("Expected client %q to be added", "John")
	}
}

func TestComputerClubManager_HandleWaiting_QueueFull(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 1, 100)
	manager.SetClub(club)

	comingEvent1 := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "Mike")
	manager.handleClientComing(comingEvent1)
	tableEvent := domain.NewTableEvent(2, startTime.Add(2*time.Minute), domain.Income, "Mike", 1)
	manager.handleTakingTable(tableEvent)

	comingEvent2 := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.handleClientComing(comingEvent2)
	event := domain.NewClientEvent(3, startTime.Add(2*time.Minute), domain.Income, "John")
	manager.handleWaiting(event)

	comingEvent3 := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "Jane")
	manager.handleClientComing(comingEvent3)
	event2 := domain.NewClientEvent(3, startTime.Add(2*time.Minute), domain.Income, "Jane")
	manager.handleWaiting(event2)

	for _, client := range manager.club.Queue {
		if client.Name == "Jane" {
			t.Errorf("Expected Jane not to be added to the queue")
		}
	}
	if manager.club.Clients["Jane"].Status != domain.Inactive {
		t.Errorf("Expected Jane's status to be inactive")
	}
	if manager.club.Events[0].GetID() != 11 {
		t.Errorf("Expected event ID to be 11")
	}
}

func TestComputerClubManager_HandleTakingTable(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	comingEvent := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.HandleEvent(comingEvent)

	tableEvent := domain.NewTableEvent(2, startTime.Add(2*time.Minute), domain.Income, "John", 1)
	manager.HandleEvent(tableEvent)

	if !manager.club.Tables[0].IsBusy() {
		t.Errorf("Expected table 1 to be busy")
	}
}

func TestComputerClubManager_HandleWaiting(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 1, 100)
	manager.SetClub(club)

	comingEvent1 := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "Mike")
	manager.HandleEvent(comingEvent1)
	tableEvent := domain.NewTableEvent(2, startTime.Add(2*time.Minute), domain.Income, "Mike", 1)
	manager.HandleEvent(tableEvent)

	comingEvent2 := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.HandleEvent(comingEvent2)
	event := domain.NewClientEvent(3, startTime.Add(2*time.Minute), domain.Income, "John")
	manager.HandleEvent(event)

	f := false
	for _, client := range manager.club.Queue {
		if client.Name == "John" {
			f = true
			break
		}
	}

	if !f {
		t.Errorf("Expected John to be added to the queue")
	}
}

func TestComputerClubManager_HandleLeaving(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	event := domain.NewClientEvent(1, startTime, domain.Income, "John")
	manager.HandleEvent(event)

	event = domain.NewClientEvent(4, startTime, domain.Income, "John")
	manager.HandleEvent(event)

	if _, exists := manager.club.Clients["John"]; exists {
		t.Errorf("Expected client %q to be removed", "John")
	}
}

func TestComputerClubManager_HandleLeaving_ClientNotInClub(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	event := domain.NewClientEvent(4, startTime, domain.Income, "John")
	manager.handleLeaving(event)

	if _, exists := manager.club.Clients["John"]; exists {
		t.Errorf("Expected client %q to not be in the club", "John")
	}
	if manager.club.Events[0].GetID() != 13 {
		t.Errorf("Expected event ID to be 13")
	}
}

func TestComputerClubManager_EndDay(t *testing.T) {
	mockOutput := output.MockOutputStrategy{}
	manager := NewComputerClubManager(&mockOutput)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 2, 100)
	manager.SetClub(club)

	newClientEvent := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.HandleEvent(newClientEvent)
	newClientEvent = domain.NewClientEvent(1, startTime.Add(2*time.Minute), domain.Income, "Mike")
	manager.HandleEvent(newClientEvent)
	newClientEvent = domain.NewClientEvent(1, startTime.Add(3*time.Minute), domain.Income, "Alex")
	manager.HandleEvent(newClientEvent)

	tableEvent := domain.NewTableEvent(2, startTime.Add(10*time.Minute), domain.Income, "John", 1)
	manager.HandleEvent(tableEvent)
	tableEvent = domain.NewTableEvent(2, startTime.Add(1*time.Hour), domain.Income, "Mike", 2)
	manager.HandleEvent(tableEvent)

	leaveEvent := domain.NewClientEvent(4, startTime.Add(2*time.Hour), domain.Income, "John")
	manager.HandleEvent(leaveEvent)

	manager.EndDay()

	if len(manager.club.Clients) != 0 {
		t.Errorf("Expected no clients remaining in the club")
	}
	if len(manager.club.Queue) != 0 {
		t.Errorf("Expected empty queue")
	}
	if len(manager.club.Events) != 0 {
		t.Errorf("Expected no events remaining")
	}
	if mockOutput.PrintRevenueCalls[0][0].Revenue() != 200 {
		t.Errorf("Expected revenue to be 200")
	}
	if mockOutput.PrintRevenueCalls[0][1].Revenue() != 700 {
		t.Errorf("Expected expenses to be 700")
	}
}

func TestComputerClubManager_HandleWaiting_ClientAlreadyHasTable(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 2, 100)
	manager.SetClub(club)

	comingEvent := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.handleClientComing(comingEvent)
	tableEvent := domain.NewTableEvent(2, startTime.Add(2*time.Minute), domain.Income, "John", 1)
	manager.handleTakingTable(tableEvent)

	waitingEvent := domain.NewClientEvent(3, startTime.Add(3*time.Minute), domain.Income, "John")
	manager.handleWaiting(waitingEvent)

	if manager.club.Clients["John"].Status != domain.Active {
		t.Errorf("Expected John's status to be active")
	}
	if manager.club.Events[0].GetID() != 13 {
		t.Errorf("Expected event ID to be 13")
	}
}

func TestComputerClubManager_HandleLeaving_ClientWasWaiting(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 1, 100)
	manager.SetClub(club)

	comingEvent := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")
	manager.handleClientComing(comingEvent)
	waitingEvent := domain.NewClientEvent(3, startTime.Add(2*time.Minute), domain.Income, "John")
	manager.handleWaiting(waitingEvent)

	leavingEvent := domain.NewClientEvent(4, startTime.Add(3*time.Minute), domain.Income, "John")
	manager.handleLeaving(leavingEvent)

	if _, exists := manager.club.Clients["John"]; exists {
		t.Errorf("Expected client %q to not be in the club", "John")
	}
	if len(manager.club.Queue) != 0 {
		t.Errorf("Expected empty queue")
	}
}
