package main

import (
	. "exercise_3/dataenums"
	"exercise_3/driver"
	"exercise_3/hwelevio"
	"exercise_3/lights"
)


func main() {
	hwelevio.Init(Addr)

	var (
		fromDriverToLight     = make(chan FromDriverToLight, 100)
	)


	go driver.ElevatorDriver(
		fromDriverToLight,
	)

	go lights.LightsHandler(
		fromDriverToLight,
	)

	select {}

}