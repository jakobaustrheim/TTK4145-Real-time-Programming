package driver

import (
	"exercise_3/dataenums"
	. "exercise_3/dataenums"
	"exercise_3/driver/timer"
	"exercise_3/hwelevio"
	"fmt"
)

func ElevatorDriver(
	payloadToLights chan<- FromDriverToLight,
) {
	var (
		floorChannel       = make(chan int)
		obstructionChannel = make(chan bool)
		doorOpenChan       = make(chan bool)
		doorClosedChan     = make(chan bool)
		motorActiveChan    = make(chan bool)
		motorInactiveChan  = make(chan bool)
		newOrderChannel   = make(chan [NFloors][NButtons]bool)
		clearedRequests    = [NFloors][NButtons]bool{}
		obstruction        bool
		payload = [NFloors][NButtons]bool{}
	)

	go hwelevio.PollFloorSensor(floorChannel)
	go hwelevio.PollObstructionSwitch(obstructionChannel)
	go timer.Timer(doorOpenChan, motorActiveChan, doorClosedChan, motorInactiveChan)
	go hwelevio.PollButtons()
	go buttonPressed(payload, dataenums.ButtonEvent )

	elevator := initelevator()
	hwelevio.SetMotorDirection(elevator.Dirn)


	payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false, Orders: elevator.Requests}

	for {
		select {
		case elevator.CurrentFloor = <-floorChannel:
			elevator.ActiveSatus = true
			motorActiveChan <- true
			switch {
			case elevator.Requests[elevator.CurrentFloor][BCab]:
				fmt.Println("")
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				doorOpenChan <- true

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case elevator.Dirn == MDUp && requestsAbove(elevator):
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			case elevator.Dirn == MDDown && requestsBelow(elevator):
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				hwelevio.SetMotorDirection(MDStop)
				elevator.CurrentBehaviour = EBDoorOpen
				motorActiveChan <- false
				doorOpenChan <- true
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

			default:
				elevator = chooseDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Dirn)
				payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
			}

		case <-doorClosedChan:
			if obstruction {
				elevator.ActiveSatus = !obstruction
				fmt.Println(!obstruction)
				doorOpenChan <- true
				continue
			}

			switch {
			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				clearedRequests[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BCab] && requestsAbove(elevator):

			case elevator.Dirn == MDUp && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				clearedRequests[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallDown]:
				clearedRequests[elevator.CurrentFloor][BHallDown] = true
				elevator.Requests[elevator.CurrentFloor][BHallDown] = false

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BCab] && requestsBelow(elevator):

			case elevator.Dirn == MDDown && elevator.Requests[elevator.CurrentFloor][BHallUp]:
				clearedRequests[elevator.CurrentFloor][BHallUp] = true
				elevator.Requests[elevator.CurrentFloor][BHallUp] = false

			//case elevator.Requests[elevator.CurrentFloor][BCab]:
			// This case was not necessary after changing chooseDirection
			// but this can have induced other errors. I have not tried yet.

			default:
				elevator = chooseDirection(elevator)
				hwelevio.SetMotorDirection(elevator.Dirn)
			}

			if elevator.Requests[elevator.CurrentFloor][BCab] {
				clearedRequests[elevator.CurrentFloor][BCab] = true
				elevator.Requests[elevator.CurrentFloor][BCab] = false
			}

			elevator.CurrentBehaviour = EBIdle

			payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: false}
			 // Reset clearedRequests to all false
			clearedRequests = [NFloors][NButtons]bool{}

		case <-motorInactiveChan:

			if elevator.CurrentBehaviour == EBMoving {
				elevator.ActiveSatus = false
			}

		case obstruction = <-obstructionChannel:
			if elevator.CurrentBehaviour == EBDoorOpen {
				elevator.ActiveSatus = !obstruction
				doorOpenChan <- !obstruction
			}

		case elevator.Requests = <-newOrderChannel:
			ElevatorPrint(elevator)
			// TODO JAKOB + ALEX 
			// CHECK THE LOGIC HERE UNDER FAT
			// THis perhaps make messy code 
			// Se if we can emulate 
			switch elevator.CurrentBehaviour {
			case EBIdle:
				switch {
				case elevator.Dirn == MDUp && (elevator.Requests[elevator.CurrentFloor][BHallUp] || elevator.Requests[elevator.CurrentFloor][BCab]):
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDUp
					doorOpenChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				
				case elevator.Dirn == MDUp && requestsAbove(elevator):
					elevator.CurrentBehaviour = EBMoving
					elevator.Dirn = MDUp
					hwelevio.SetMotorDirection(elevator.Dirn)
					motorActiveChan <- true

				case elevator.Dirn == MDDown && (elevator.Requests[elevator.CurrentFloor][BHallDown] || elevator.Requests[elevator.CurrentFloor][BCab]):
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDDown
					doorOpenChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				
				case  elevator.Dirn == MDDown && requestsBelow(elevator):
					elevator.CurrentBehaviour = EBMoving
					elevator.Dirn = MDDown
					hwelevio.SetMotorDirection(elevator.Dirn)
					motorActiveChan <- true
				
				case  elevator.Requests[elevator.CurrentFloor][BHallUp]:
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDUp
					doorOpenChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case  elevator.Requests[elevator.CurrentFloor][BHallDown]:
					elevator.CurrentBehaviour = EBDoorOpen
					elevator.Dirn = MDDown
					doorOpenChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}

				case elevator.Requests[elevator.CurrentFloor][BCab]:
					elevator.CurrentBehaviour = EBDoorOpen
					doorOpenChan <- true
					payloadToLights <- FromDriverToLight{CurrentFloor: elevator.CurrentFloor, DoorLight: true}
				
				default:	
					elevator = chooseDirection(elevator)
					hwelevio.SetMotorDirection(elevator.Dirn)
				}

			case EBMoving:
			case EBDoorOpen:
			}
		}
	}
}