package main

import (
	"Driver-go/elevio"
)

func RequestsAbove(elevator *Elevator) bool { // skal returnere true/false om det er noen aktive orders i etasjer over
	for f := elevator.CurrentFloor + 1; f < len(elevator.Orders); f++ {
		for _, order := range elevator.Orders[f] {
			if order {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(elevator *Elevator) bool { // skal returnere true/false om det er noen aktive orders i etasjer under
	for f := elevator.CurrentFloor - 1; f < len(elevator.Orders); f-- {
		for _, order := range elevator.Orders[f] {
			if order {
				return true
			}
		}
	}
	return false
}

func RequestsHere(elevator *Elevator) bool {
	for _, order := range elevator.Orders[elevator.CurrentFloor] {
		if order {
			return true
		}
	}
	return false
}

func AddOrder(elevator *Elevator, floor int, button elevio.ButtonType) {
	elevator.Orders[floor][button] = true
	elevio.SetButtonLamp(button, floor, true)

}

func ClearRequestsAtFloor(elevator *Elevator) {
	// iterer over knappetypene (0, 1, 2)
	for b := 0; b < 3; b++ {
		// sletter alle bestillinger i den etasjen man er i:
		elevator.Orders[elevator.CurrentFloor][b] = false
		elevio.SetButtonLamp(elevio.ButtonType(b), elevator.CurrentFloor, false)
	}
}

func ClearAllRequests(elevator *Elevator) {
	for f := 0; f < len(elevator.Orders); f++ {
		for b := 0; b < 3; b++ {
			elevator.Orders[f][b] = false
			elevio.SetButtonLamp(elevio.ButtonType(b), f, false)
		}
	}
}

func ChooseDirection(elevator *Elevator) elevio.MotorDirection { //velger retning basert på nåværende retning og bestillinger
	switch elevator.Direction {
	case elevio.MD_Up:
		if RequestsAbove(elevator) {
			return elevio.MD_Up
		}
		if RequestsHere(elevator) || RequestsBelow(elevator) {
			return elevio.MD_Down
		}
	case elevio.MD_Down:
		if RequestsBelow(elevator) {
			return elevio.MD_Down
		}
		if RequestsHere(elevator) || RequestsAbove(elevator) {
			return elevio.MD_Up
		}
	case elevio.MD_Stop: // dersom den står stille prioriterer den bestillinger som er over
		if RequestsAbove(elevator) {
			return elevio.MD_Up
		}
		if RequestsBelow(elevator) {
			return elevio.MD_Down
		}
	}
	return elevio.MD_Stop

}

func ShouldStop(elevator *Elevator) bool {
	switch elevator.Direction {
	case elevio.MD_Up:
		return RequestsHere(elevator) || !RequestsAbove(elevator)
	case elevio.MD_Down:
		return RequestsHere(elevator) || !RequestsBelow(elevator)
	case elevio.MD_Stop:
		return RequestsHere(elevator)
	}
	return false
}
