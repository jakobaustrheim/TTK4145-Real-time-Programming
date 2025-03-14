package lights

import (
	. "exercise_3/dataenums"
	"exercise_3/hwelevio"
)

func LightsHandler(
	payloadFromDriver <-chan FromDriverToLight,
) {
	for {
		select {
		case payload := <-payloadFromDriver:
			hwelevio.SetFloorIndicator(payload.CurrentFloor)
			hwelevio.SetDoorOpenLamp(payload.DoorLight)
			for floor := 0; floor < NFloors; floor++ {
				for button := 0; button < NButtons; button++ {
					hwelevio.SetButtonLamp(Button(button), floor, Order[floor][button] == OrderAssigned)
				}
			}

		}
	}
}
