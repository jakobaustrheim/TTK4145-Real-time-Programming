// Use `go run foo.go` to run your program

package main

import (
    . "fmt"
    "runtime"
    // "time"
)

var i = 0

// func incrementing() {
//     //TODO: increment i 1000000 times
//     for j := 0; j < 1000000; j++ {
//         i++
//     }
// }

// func decrementing() {
//     //TODO: decrement i 1000000 times
//     for j := 0; j < 1000000; j++ {
//         i--
//     }
// }

// func main() {
//     // What does GOMAXPROCS do? What happens if you set it to 1?
//     runtime.GOMAXPROCS(2)  
    
	
//     // TODO: Spawn both functions as goroutines

//     go incrementing()
//     go decrementing()

	
//     // We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
//     // We will do it properly with channels soon. For now: Sleep.

//     time.Sleep(500*time.Millisecond)
//     Println("The magic number is:", i)
// }

type message struct {
    action string
    done   chan bool
}

func main() {

    runtime.GOMAXPROCS(2)


    incrementChan := make(chan message)
    decrementChan := make(chan message)
    doneChan := make(chan bool)

    go func() {
        for {
            select {
            case msg := <-incrementChan:
                i++
                msg.done <- true
            case msg := <-decrementChan:
                i--
                msg.done <- true
            }
        }
    }()

    go func() {
        for j := 0; j < 1000000; j++ {
            done := make(chan bool)
            incrementChan <- message{action: "increment", done: done}
            <-done
        }
        doneChan <- true
    }()

    go func() {
        for j := 0; j < 1000000-1; j++ {
            done := make(chan bool)
            decrementChan <- message{action: "decrement", done: done}
            <-done
        }
        doneChan <- true
    }()

    <-doneChan
    <-doneChan

    Println("The magic number is:", i)
}