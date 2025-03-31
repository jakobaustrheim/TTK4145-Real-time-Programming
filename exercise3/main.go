package main

import (
	"Exercise3/elevatorlogic"
	"Exercise3/elevio"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	// var d elevio.MotorDirection = elevio.MD_Up
	// //elevio.SetMotorDirection(d)

	elevatorlogic.InitElevatorLogic()

	select {}

}
