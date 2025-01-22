package main

import "Driver-go/elevio"

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

// func initElevator() Elevator {
// 	return Elevator{
// 		CurrentFloor : -1,
// 		Direction : elevio.MD_Stop,
// 		State : Idle,
// 	}
// }


