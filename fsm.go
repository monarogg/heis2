package main

import (
	"Driver-go/elevio"
	"time"
)

type ElevatorState int

const (
	Idle ElevatorState = 0
	Moving ElevatorState = 1
	DoorOpen ElevatorState= 2
)

type Elevator struct {
	CurrentFloor int
	Direction elevio.MotorDirection
	State ElevatorState
	Orders [4][3]bool
}

type ElevatorConfig struct {
	DoorOpenConfiguration time.Duration
}

// func initElevator() Elevator {
// 	return Elevator{
// 		CurrentFloor : -1,
// 		Direction : elevio.MD_Stop,
// 		State : Idle,
// 	}
// }

func initializeFSM() Elevator {		// funksjonen returnerer ferdiginitialisert instans av strukturen Elevator
	
	elevator := Elevator {
		CurrentFloor: 0,	//starter i første etasje
		Direction: elevio.MD_Stop,	// motoren skal stå i ro
		State: Idle,	//starter som inaktiv
		Orders: [4][3]bool{},		//ingen bestillinger
	}

	return elevator
}

func OnRequestButtonPress(elevator *Elevator, button elevio.ButtonType) {
	elevator.Orders[button.Floot][button.Button] = true

	switch elevator.State {

	case DoorOpen:

	case Moving:

	case Idle:

	}
}

func OnFloorArrival(elevator *Elevator, floor int) {
	elevator.CurrentFloor = floor	// oppdaterer current floor
	elevio.SetFloorIndicator(floor)

	if ShouldStop(elevator) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true) // skal door open lamp være på her?
		ClearAllRequests(elevator)
		elevator.State = Idle // usikker på om state skal være Idle??

		// starte dørens tidsavbrudd????

	}
}
// funksjon som skal brukes når døren har vært åpen tilstrekkelig lenge (dør skal lukkes osv.):
func OnDoorTimeout(elevator *Elevator) {
	elevio.SetDoorOpenLamp(false)

	elevator.Direction = ChooseDirection(elevator)
	if elevator.Direction == elevio.MD_Stop {
		elevator.State = Idle
	} else {
		elevio.SetMotorDirection(elevator.Direction)
		elevator.State = Moving
	}
}

func OnStopButtonPress(elevator *Elevator) {

	elevio.SetStopLamp(true)
	elevio.SetMotorDirection((elevio.MD_Stop))
	ClearAllRequests(elevator)
	elevator.State = Idle
	UpdateLights(elevator)

}

func ChooseDirection(elevator *Elevator) elevio.MotorDirection {
	// iterer over alle etasjene over current floor:
	for f := elevator.CurrentFloor + 1;  f < len(elevator.Orders); f++ {
		// iterer over alle knappetrykk for etasjen f:
		for _; order := range elevator.Orders[f] {
			// sjekker om det er en order (hvis den er true):
			if order {
				// skal da bevege seg oppover
				return elevio.MD_Up
			}
		}
	}

	for f := elevator.CurrentFloor - 1; f >= 0; f-- {
		for _; order := range elevator.Orders[f] {
			if order {
				return elevio.MD_Down
			}
		}
	}

	// dersom ikke går inn i noen av if setningene, skal heisen stoppe:
	return elevio.MD_Stop
}

func ClearRequestsAtFloor(elevator *Elevator) {
	// iterer over knappetypene (0, 1, 2)
	for b := 0; b < 3; b++ {
		// sletter alle bestillinger i den etasjen man er i:
		elevator.Orders[elevator.CurrentFloor][b] = false
		elevio.SetButtonLamp(elevio.ButtonType(b), elevator.CurrentFloor, false)
	}
}

func UpdateLights(elevator *Elevator) {
	// skal oppdatere button lights basert på aktive orders i Orders
	for f := 0; f < len(elevator.Orders); f++ {
		for b := 0; b < 3; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elevator.Orders[f][b])
		}
	}
}

func ShouldStop(elevator * Elevator) bool {
	return elevator.Orders[elevator.CurrentFloor][elevio.BT_HallUp] || 
		elevator.Orders[elevator.CurrentFloor][elevio.BT_HallDown] ||
		elevator.Orders[elevator.CurrentFloor][elevio.BT_Cab]
}

func ClearAllRequests(elevator *Elevator) {
	for f := 0; f < len(elevator.Orders); f++ {
		for b := 0; b < 3; b++ {
			elevator.Orders[f][b] = false
			elevio.SetButtonLamp(elevio.ButtonType(b), f, false)
		}
	}
}