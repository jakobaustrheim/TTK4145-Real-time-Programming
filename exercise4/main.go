package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

const (
	sendAddr       = "255.255.255.255:20026"
	receiveAddr    = ":20026"
	heartbeatMsg   = "heartbeat"
	heartbeatSleep = 500
)

// Function to start a backup process that will become primary if needed.
func startBackupProcess() {
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
}

// The primary process sends heartbeats to the backup.
func primaryProcess(count int) {
	sendUDPAddr, _ := net.ResolveUDPAddr("udp", sendAddr)
	conn, _ := net.DialUDP("udp", nil, sendUDPAddr)
	defer conn.Close()
	for {
		msg := strconv.Itoa(count)
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Primary failed to send heartbeat:", err)
			return
		}
		fmt.Printf("%d \n", count)
		count++
		time.Sleep(heartbeatSleep * time.Millisecond)
	}
}

// The backup process listens for heartbeats from the primary.
func backupProcess() {
	count := 1
	fmt.Printf("The place for counting. Count on us\n")
	receiveUDPAddr, _ := net.ResolveUDPAddr("udp", receiveAddr)
	conn, _ := net.ListenUDP("udp", receiveUDPAddr)
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 2 * time.Millisecond))
		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				conn.Close()
				startBackupProcess()
				primaryProcess(count)
				return
			} else {
				fmt.Println("Error reading from UDP:", err)
				return
			}
		}

		msg := string(buffer[:n])
		recievedCount, _ := strconv.Atoi(msg)
		if recievedCount != 0 {
			count = recievedCount
		}
	}
}

func main() {
	backupProcess()
}
