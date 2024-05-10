package service

import (
	"gomputerClub/internal/adapters/output"
	"gomputerClub/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestComputerClubManager_HandleClientComing_ClubClosed(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	event := domain.NewClientEvent(1, startTime.Add(-1*time.Minute), domain.Income, "John")
	manager.handleClientComing(event)
	assert.NotContains(t, manager.club.Clients, "John")

	event = domain.NewClientEvent(1, startTime.Add(9*time.Hour), domain.Income, "John")
	manager.handleClientComing(event)
	assert.NotContains(t, manager.club.Clients, "John")
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

	assert.Equal(t, 1, manager.club.Clients["John"].TableNumber)
	assert.Equal(t, domain.Inactive, manager.club.Clients["Mike"].Status)
	assert.True(t, manager.club.Tables[0].IsBusy())
}

func TestComputerClubManager_HandleClientComing(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	event := domain.NewClientEvent(1, startTime.Add(1*time.Minute), domain.Income, "John")

	manager.HandleEvent(event)

	assert.Contains(t, manager.club.Clients, "John")
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

	assert.NotContains(t, manager.club.Queue, manager.club.Clients["Jane"])
	assert.Equal(t, domain.Inactive, manager.club.Clients["Jane"].Status)
	assert.Equal(t, 11, manager.club.Events[0].GetID())
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

	assert.True(t, manager.club.Tables[0].IsBusy())
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

	assert.Contains(t, manager.club.Queue, manager.club.Clients["John"])
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

	assert.NotContains(t, manager.club.Clients, "John")
}

func TestComputerClubManager_HandleLeaving_ClientNotInClub(t *testing.T) {
	manager := NewComputerClubManager(nil)
	startTime, _ := time.Parse("15:04", "10:00")
	club := domain.NewComputerClub(startTime, startTime.Add(8*time.Hour), 10, 100)
	manager.SetClub(club)

	event := domain.NewClientEvent(4, startTime, domain.Income, "John")
	manager.handleLeaving(event)

	assert.NotContains(t, manager.club.Clients, "John")
	assert.Equal(t, 13, manager.club.Events[0].GetID())
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

	assert.Empty(t, manager.club.Clients)
	assert.Empty(t, manager.club.Queue)
	assert.Empty(t, manager.club.Events)

	assert.Equal(t, 200, mockOutput.PrintRevenueCalls[0][0].Revenue())
	assert.Equal(t, 700, mockOutput.PrintRevenueCalls[0][1].Revenue())
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

	assert.Equal(t, domain.Active, manager.club.Clients["John"].Status)
	assert.Equal(t, 13, manager.club.Events[0].GetID())
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

	assert.NotContains(t, manager.club.Clients, "John")
	assert.Empty(t, manager.club.Queue)
}
