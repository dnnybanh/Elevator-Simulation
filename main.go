package main

import (
	"fmt"

	"main.go/logger" // Update this import path based on your project's structure.
)

type Direction int

const (
	UP Direction = iota
	DOWN
	IDLE
)

type Request struct {
	PickupFloor  int
	DropoffFloor int
}

type Elevator struct {
	headPosition int
	upQueue      []Request
	downQueue    []Request
	passengers   []Request
	direction    Direction
}

func NewElevator(headPosition int) *Elevator {
	e := &Elevator{
		headPosition: headPosition,
		direction:    IDLE,
		upQueue:      []Request{},
		downQueue:    []Request{},
		passengers:   []Request{},
	}
	logger.Info(fmt.Sprintf("Elevator initialized at floor %d", e.headPosition))
	return e
}

func (e *Elevator) addRequest(pickup, dropoff int) {
	request := Request{PickupFloor: pickup, DropoffFloor: dropoff}
	if pickup >= e.headPosition {
		e.upQueue = append(e.upQueue, request)
	} else {
		e.downQueue = append(e.downQueue, request)
	}
}

func (e *Elevator) sortQueues() {
	// Sorting upQueue in ascending order based on PickupFloor
	for i := 0; i < len(e.upQueue)-1; i++ {
		for j := 0; j < len(e.upQueue)-i-1; j++ {
			if e.upQueue[j].PickupFloor > e.upQueue[j+1].PickupFloor {
				e.upQueue[j], e.upQueue[j+1] = e.upQueue[j+1], e.upQueue[j]
			}
		}
	}

	// Sorting downQueue in descending order based on PickupFloor
	// This ensures that as the elevator moves down, it services the highest floors first.
	for i := 0; i < len(e.downQueue)-1; i++ {
		for j := 0; j < len(e.downQueue)-i-1; j++ {
			if e.downQueue[j].PickupFloor < e.downQueue[j+1].PickupFloor {
				e.downQueue[j], e.downQueue[j+1] = e.downQueue[j+1], e.downQueue[j]
			}
		}
	}
}

func (e *Elevator) processLOOK() {
	logger.Info("Starting to process requests using LOOK algorithm.")

	// Process the upQueue first
	e.direction = UP
	for len(e.upQueue) > 0 {
		nextRequest := e.upQueue[0]
		e.upQueue = e.upQueue[1:]
		e.headPosition = nextRequest.PickupFloor
		// Logic to simulate picking up the passenger
		logger.Info(fmt.Sprintf("Elevator moving UP to pick up passenger at floor %d for drop-off at floor %d.", nextRequest.PickupFloor, nextRequest.DropoffFloor))
		// Add passenger to the onboard list for drop-off
		e.passengers = append(e.passengers, nextRequest)
		// Check and drop off any passengers if their drop-off floor is reached
		e.checkAndDropOffPassengers()
	}

	// After finishing the upQueue, check if downQueue has requests
	if len(e.downQueue) > 0 {
		e.direction = DOWN
		for len(e.downQueue) > 0 {
			nextRequest := e.downQueue[0]
			e.downQueue = e.downQueue[1:]
			e.headPosition = nextRequest.PickupFloor
			// Logic to simulate picking up the passenger
			logger.Info(fmt.Sprintf("Elevator moving DOWN to pick up passenger at floor %d for drop-off at floor %d.", nextRequest.PickupFloor, nextRequest.DropoffFloor))
			// Add passenger to the onboard list for drop-off
			e.passengers = append(e.passengers, nextRequest)
			// Check and drop off any passengers if their drop-off floor is reached
			e.checkAndDropOffPassengers()
		}
	}

	e.direction = IDLE
	logger.Info("Finished processing all requests. Elevator is idle.")
}

// checkAndDropOffPassengers checks the passengers slice for any passengers that need to be dropped off at the current floor.
func (e *Elevator) checkAndDropOffPassengers() {// Before dropping off, log the attempt
	logger.Info(fmt.Sprintf("Checking for drop-offs at floor %d", e.headPosition))

    for i := 0; i < len(e.passengers); {
        if e.passengers[i].DropoffFloor == e.headPosition {
			logger.Info(fmt.Sprintf("Dropped off passenger at floor %d.", e.headPosition))
			e.logPassengerState()
            // Remove the passenger from the slice by replacing it with the last element and truncating the slice.
            e.passengers[i] = e.passengers[len(e.passengers)-1]
            e.passengers = e.passengers[:len(e.passengers)-1]
            // Don't increment i, as we need to check the new passenger that was swapped into this position.
        } else {
            // Only increment i if no passenger was dropped off, as we didn't swap any passengers into this position.
            i++
        }
    }
}

func (e *Elevator) logCurrentState() {
	header := []string{"Queue", "Pickup Floor", "Dropoff Floor"}
	var rows [][]string

	for _, req := range e.upQueue {
		rows = append(rows, []string{"UP", fmt.Sprintf("%d", req.PickupFloor), fmt.Sprintf("%d", req.DropoffFloor)})
	}
	for _, req := range e.downQueue {
		rows = append(rows, []string{"DOWN", fmt.Sprintf("%d", req.PickupFloor), fmt.Sprintf("%d", req.DropoffFloor)})
	}
	logger.InfoTable(header, rows)
}

func (e *Elevator) logPassengerState() {
    if len(e.passengers) == 0 {
        logger.Info("No passengers currently onboard.")
        return
    }

    header := []string{"Onboard Passenger", "Dropoff Floor"}
    var rows [][]string

    for i, passenger := range e.passengers {
        row := []string{
            fmt.Sprintf("Passenger %d", i+1),
            fmt.Sprintf("%d", passenger.DropoffFloor),
        }
        rows = append(rows, row)
    }

    logger.InfoTable(header, rows)
}

func main() {
    elevator := NewElevator(5) // Starting floor for the elevator

    // Adding requests with both pickup and drop-off floors
    elevator.addRequest(3, 6)  // Pickup at floor 3, drop off at floor 6
    elevator.addRequest(7, 2)  // Pickup at floor 7, drop off at floor 2
    elevator.addRequest(8, 3)  // Pickup at floor 8, drop off at floor 3
    elevator.addRequest(2, 8)  // Pickup at floor 2, drop off at floor 8
    elevator.addRequest(12, 1) // Pickup at floor 12, drop off at floor 1
    elevator.addRequest(1, 10) // Pickup at floor 1, drop off at floor 10
	elevator.logCurrentState()

    elevator.sortQueues()

    // Process requests using the updated LOOK algorithm with passenger exchanges
    elevator.processLOOK()
}

