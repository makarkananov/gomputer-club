package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTable_StartSession_SetsIsBusyAndLastStart(t *testing.T) {
	table := NewTable()
	startTime := time.Now()

	table.StartSession(startTime)

	assert.True(t, table.isBusy)
	assert.Equal(t, startTime, table.lastStart)
}

func TestTable_EndSession_SetsIsNotBusyAndUpdatesBusyDurationAndRevenue(t *testing.T) {
	table := NewTable()
	startTime := time.Now()
	table.StartSession(startTime)

	endTime := startTime.Add(2 * time.Hour)
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	assert.False(t, table.isBusy)
	assert.Equal(t, 2*time.Hour, table.busyDuration)
	assert.Equal(t, 200, table.revenue)
}

func TestTable_EndSession_WhenTableIsNotBusy_DoesNotChangeState(t *testing.T) {
	table := NewTable()
	endTime := time.Now()
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	assert.False(t, table.isBusy)
	assert.Equal(t, 0*time.Second, table.busyDuration)
	assert.Equal(t, 0, table.revenue)
}

func TestTable_EndSession_WhenSessionIsLessThanOneHour_ChargesForOneHour(t *testing.T) {
	table := NewTable()
	startTime := time.Now()
	table.StartSession(startTime)

	endTime := startTime.Add(30 * time.Minute)
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	assert.False(t, table.isBusy)
	assert.Equal(t, 30*time.Minute, table.busyDuration)
	assert.Equal(t, 100, table.revenue)
}

func TestTable_EndSession_WhenSessionIsMoreThanOneHour_ChargesForActualHours(t *testing.T) {
	table := NewTable()
	startTime := time.Now()
	table.StartSession(startTime)

	endTime := startTime.Add(90 * time.Minute)
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	assert.False(t, table.isBusy)
	assert.Equal(t, 90*time.Minute, table.busyDuration)
	assert.Equal(t, 200, table.revenue)
}
