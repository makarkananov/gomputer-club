package input

import (
	"bufio"
	"fmt"
	"gomputerClub/internal/core/domain"
	"gomputerClub/internal/core/port"
	"os"
	"strconv"
	"strings"
	"time"
)

// FileController is a controller for reading and running events from a file.
type FileController struct {
	fileName string
	ccm      port.ComputerClubManager
}

// NewFileController creates a new FileController with the given file name and ComputerClubManager.
func NewFileController(fileName string, ccm port.ComputerClubManager) *FileController {
	return &FileController{
		fileName: fileName,
		ccm:      ccm,
	}
}

// Read reads all events from the file and returns them.
func (fic *FileController) Read() ([]domain.Event, error) {
	file, err := os.Open(fic.fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var numberOfTables int
	if scanner.Scan() {
		numberOfTables, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("error parsing number of tables: %w", err)
		}
	}

	var openTime, closeTime time.Time
	if scanner.Scan() {
		times := strings.Split(scanner.Text(), " ")
		openTime, err = time.Parse("15:04", times[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing open time: %w", err)
		}

		closeTime, err = time.Parse("15:04", times[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing close time: %w", err)
		}
	}

	var costPerHour int
	if scanner.Scan() {
		costPerHour, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("error parsing cost per hour: %w", err)
		}
	}

	fic.ccm.SetClub(domain.NewComputerClub(openTime, closeTime, numberOfTables, costPerHour))

	events := make([]domain.Event, 0)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		event, err := fic.parseEvent(parts)
		if err != nil {
			return nil, fmt.Errorf("error parsing event: %w", err)
		}

		events = append(events, event)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return events, nil
}

// parseEvent parses an event from the given components.
func (fic *FileController) parseEvent(components []string) (domain.Event, error) {
	if len(components) < 3 {
		return nil, fmt.Errorf("not enough components in line: %s", strings.Join(components, " "))
	}

	eventTime, err := time.Parse("15:04", components[0])
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing event time in line: %s, error: %w",
			strings.Join(components, " "),
			err,
		)
	}

	eventID, err := strconv.Atoi(components[1])
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing event ID in line: %s, error: %w",
			strings.Join(components, " "),
			err,
		)
	}

	clientName := components[2]

	clientEvent := domain.ClientEvent{
		BaseEvent: domain.BaseEvent{
			ID:        eventID,
			Time:      eventTime,
			EventType: domain.Income,
		},
		ClientName: clientName,
	}

	if eventID == 2 {
		if len(components) != 4 {
			return nil, fmt.Errorf(
				"not enough components for event type 2 in line: %s",
				strings.Join(components, " "),
			)
		}

		tableNumber, err := strconv.Atoi(components[3])
		if err != nil {
			return nil, fmt.Errorf(
				"error parsing table number in line: %s, error: %w",
				strings.Join(components, " "),
				err,
			)
		}

		return domain.TableEvent{
			ClientEvent: clientEvent,
			TableNumber: tableNumber,
		}, nil

	}

	if len(components) > 3 {
		return nil, fmt.Errorf("too many components in line: %s", strings.Join(components, " "))
	}

	return clientEvent, nil
}

// Run runs the given events on the ComputerClubManager.
func (fic *FileController) Run(events []domain.Event) {
	for _, event := range events {
		fic.ccm.HandleEvent(event)
	}

	fic.ccm.EndDay()
}
