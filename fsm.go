package main

import (
	"Driver-go/elevio"
	"time"
)

type ElevatorState int

const (
	Idle     ElevatorState = 0
	Moving   ElevatorState = 1
	DoorOpen ElevatorState = 2
)

type Elevator struct {
	CurrentFloor int
	Direction    elevio.MotorDirection
	State        ElevatorState
	Orders       [4][3]bool
	Config       ElevatorConfig
}

type ElevatorConfig struct {
	DoorOpenDuration time.Duration
}

// func initElevator() Elevator {
// 	return Elevator{
// 		CurrentFloor : -1,
// 		Direction : elevio.MD_Stop,
// 		State : Idle,
// 	}
// }

func initializeFSM() Elevator { // funksjonen returnerer ferdiginitialisert instans av strukturen Elevator

	elevator := Elevator{
		CurrentFloor: 0,              //starter i første etasje
		Direction:    elevio.MD_Stop, // motoren skal stå i ro
		State:        Idle,           //starter som inaktiv
		Orders:       [4][3]bool{},   //ingen bestillinger
	}

	return elevator
}

func OnRequestButtonPress(elevator *Elevator, btnFloor int, btnType elevio.ButtonType) {
	switch elevator.State {
	case DoorOpen:
		if elevator.CurrentFloor == btnFloor {
			StartDoorTimer(elevator, elevator.Config.DoorOpenDuration)
			ClearRequestsAtFloor(elevator)
		}
	case Moving:
	}
}

func StartDoorTimer(elevator *Elevator, duration time.Duration) {
	time.AfterFunc(duration, func() { //time.Afterfunc starter en timer som varer i duration
		OnDoorTimeout(elevator)

	})
}

func OnFloorArrival(elevator *Elevator, floor int) {
	elevator.CurrentFloor = floor // oppdaterer current floor
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
		elevator.State = Moving
		elevio.SetMotorDirection(elevator.Direction)

	}
}

func OnStopButtonPress(elevator *Elevator) {

	elevio.SetStopLamp(true)
	elevio.SetMotorDirection((elevio.MD_Stop))
	ClearAllRequests(elevator)
	elevator.State = Idle
	UpdateLights(elevator)

}

func UpdateLights(elevator *Elevator) {
	// skal oppdatere button lights basert på aktive orders i Orders
	for f := 0; f < len(elevator.Orders); f++ {
		for b := 0; b < 3; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elevator.Orders[f][b])
		}
	}
}
