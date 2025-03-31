package elevatorlogic

import (
	"Exercise3/elevio"
	"fmt"
	"time"
)

type Orders struct {
	Up     [4]int
	Down   [4]int
	Inside [4]int
}

const (
	MovingUp = iota
	MovingDown
	GoingUpToServeDown
	GoingDownToServeUp
	Undefined
)

func setOrder(btn elevio.ButtonEvent, orders *Orders) {
	switch btn.Button {
	case elevio.BT_HallUp:
		orders.Up[btn.Floor] = 1
	case elevio.BT_HallDown:
		orders.Down[btn.Floor] = 1
	case elevio.BT_Cab:
		orders.Inside[btn.Floor] = 1
	}
	elevio.SetButtonLamp(btn.Button, btn.Floor, true)
}

func clearOrder(btn elevio.ButtonEvent, orders *Orders) {
	switch btn.Button {
	case elevio.BT_HallUp:
		orders.Up[btn.Floor] = 0
	case elevio.BT_HallDown:
		orders.Down[btn.Floor] = 0
	case elevio.BT_Cab:
		orders.Inside[btn.Floor] = 0
	}
	elevio.SetButtonLamp(btn.Button, btn.Floor, false)
}

func determineDirection(state *int, orders *Orders, currentFloor int) elevio.MotorDirection {
	if *state == MovingUp || *state == GoingUpToServeDown {
		for f := currentFloor + 1; f < 4; f++ {
			if orders.Up[f] == 1 || orders.Inside[f] == 1 {
				*state = MovingUp
				return elevio.MD_Up
			}
		}
		for f := 3; f >= 0; f-- {
			if orders.Down[f] == 1 || orders.Inside[f] == 1 {
				*state = MovingDown
				return motorDir(currentFloor, f)
			}
		}
		for f := 0; f < currentFloor; f++ {
			if orders.Up[f] == 1 {
				*state = GoingDownToServeUp
				return elevio.MD_Down
			}
		}
	} else if *state == MovingDown || *state == GoingDownToServeUp {
		for f := currentFloor - 1; f >= 0; f-- {
			if orders.Down[f] == 1 || orders.Inside[f] == 1 {
				*state = MovingDown
				return elevio.MD_Down
			}
		}
		for f := 0; f < 4; f++ {
			if orders.Up[f] == 1 || orders.Inside[f] == 1 {
				*state = MovingUp
				return motorDir(currentFloor, f)
			}
		}
		for f := 3; f > currentFloor; f-- {
			if orders.Down[f] == 1 {
				*state = GoingUpToServeDown
				return elevio.MD_Up
			}
		}
	}
	return elevio.MD_Stop
}

func motorDir(current, target int) elevio.MotorDirection {
	if current < target {
		return elevio.MD_Up
	} else if current > target {
		return elevio.MD_Down
	}
	return elevio.MD_Stop
}

func shouldServeNow(floor int, orders *Orders, state int) bool {
	if orders.Inside[floor] == 1 {
		return true
	}
	if (state == MovingUp || state == GoingDownToServeUp) && orders.Up[floor] == 1 {
		return true
	}
	if (state == MovingDown || state == GoingUpToServeDown) && orders.Down[floor] == 1 {
		return true
	}
	return false
}

func serveCurrentFloor(floor int, orders *Orders, state int, obstruction <-chan bool, ready chan<- bool) {
	clearOrder(elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}, orders)
	if state == MovingDown {
		clearOrder(elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallDown}, orders)
	} else if state == MovingUp {
		clearOrder(elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallUp}, orders)
	}
	elevio.SetDoorOpenLamp(true)
	timer := time.NewTimer(3 * time.Second)
	obstructed := false

	for {
		select {
		case obstructed = <-obstruction:
			if !obstructed {
				timer.Reset(3 * time.Second)
			}
		case <-timer.C:
			if !obstructed {
				elevio.SetDoorOpenLamp(false)
				ready <- true
				return
			}
		}
	}
}

func elevatorLoop(state int, orders Orders, floor int,
	btns <-chan elevio.ButtonEvent, floors <-chan int, obstruction <-chan bool, stop <-chan bool, ready chan bool) {

	stationary := true
	canMove := true

	for {
		select {
		case btn := <-btns:
			setOrder(btn, &orders)
			if stationary && btn.Floor == floor {
				elevio.SetMotorDirection(elevio.MD_Stop)
				canMove = false
				go serveCurrentFloor(floor, &orders, state, obstruction, ready)
			}
			if canMove {
				dir := determineDirection(&state, &orders, floor)
				elevio.SetMotorDirection(dir)
			}

		case newFloor := <-floors:
			floor = newFloor
			elevio.SetFloorIndicator(floor)
			if shouldServeNow(floor, &orders, state) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				canMove = false
				go serveCurrentFloor(floor, &orders, state, obstruction, ready)
			}

		case <-ready:
			canMove = true
			dir := determineDirection(&state, &orders, floor)
			elevio.SetMotorDirection(dir)
			stationary = dir == elevio.MD_Stop
			if shouldServeNow(floor, &orders, state) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				canMove = false
				go serveCurrentFloor(floor, &orders, state, obstruction, ready)
			}

		case <-stop:
			fmt.Println("Emergency stop activated")
		}
	}
}

func InitElevatorLogic() {
	var (
		state        = Undefined
		orders       = Orders{}
		currentFloor int

		drvButtons   = make(chan elevio.ButtonEvent)
		drvFloors    = make(chan int)
		drvObstruct  = make(chan bool)
		drvStop      = make(chan bool)
		readyToMove  = make(chan bool)
	)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstruct)
	go elevio.PollStopButton(drvStop)

	for btn := 0; btn < 3; btn++ {
		for flr := 0; flr < 4; flr++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), flr, false)
		}
	}

	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(elevio.MD_Down)
	currentFloor = <-drvFloors
	state = MovingUp
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(currentFloor)

	go elevatorLoop(state, orders, currentFloor, drvButtons, drvFloors, drvObstruct, drvStop, readyToMove)
}
