package domain

import (
	"testing"
	"time"
)

func TestTable_StartSession_SetsIsBusyAndLastStart(t *testing.T) {
	table := NewTable()
	startTime := time.Now()

	table.StartSession(startTime)

	if !table.isBusy {
		t.Errorf("Expected table to be busy")
	}
	if table.lastStart != startTime {
		t.Errorf("Expected last start time to be %v, got %v", startTime, table.lastStart)
	}
}

func TestTable_EndSession_SetsIsNotBusyAndUpdatesBusyDurationAndRevenue(t *testing.T) {
	table := NewTable()
	startTime := time.Now()
	table.StartSession(startTime)

	endTime := startTime.Add(2 * time.Hour)
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	if table.isBusy {
		t.Errorf("Expected table not to be busy")
	}
	if table.busyDuration != 2*time.Hour {
		t.Errorf("Expected busy duration to be 2 hours, got %v", table.busyDuration)
	}
	if table.revenue != 200 {
		t.Errorf("Expected revenue to be 200, got %v", table.revenue)
	}
}

func TestTable_EndSession_WhenTableIsNotBusy_DoesNotChangeState(t *testing.T) {
	table := NewTable()
	endTime := time.Now()
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	if table.isBusy {
		t.Errorf("Expected table not to be busy")
	}
	if table.busyDuration != 0*time.Second {
		t.Errorf("Expected busy duration to be 0 seconds, got %v", table.busyDuration)
	}
	if table.revenue != 0 {
		t.Errorf("Expected revenue to be 0, got %v", table.revenue)
	}
}

func TestTable_EndSession_WhenSessionIsLessThanOneHour_ChargesForOneHour(t *testing.T) {
	table := NewTable()
	startTime := time.Now()
	table.StartSession(startTime)

	endTime := startTime.Add(30 * time.Minute)
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	if table.isBusy {
		t.Errorf("Expected table not to be busy")
	}
	if table.busyDuration != 30*time.Minute {
		t.Errorf("Expected busy duration to be 30 minutes, got %v", table.busyDuration)
	}
	if table.revenue != 100 {
		t.Errorf("Expected revenue to be 100, got %v", table.revenue)
	}
}

func TestTable_EndSession_WhenSessionIsMoreThanOneHour_ChargesForActualHours(t *testing.T) {
	table := NewTable()
	startTime := time.Now()
	table.StartSession(startTime)

	endTime := startTime.Add(90 * time.Minute)
	costPerHour := 100

	table.EndSession(endTime, costPerHour)

	if table.isBusy {
		t.Errorf("Expected table not to be busy")
	}
	if table.busyDuration != 90*time.Minute {
		t.Errorf("Expected busy duration to be 90 minutes, got %v", table.busyDuration)
	}
	if table.revenue != 200 {
		t.Errorf("Expected revenue to be 200, got %v", table.revenue)
	}
}
