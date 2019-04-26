package main

import (
	"fmt"
	"os"
	"time"
)

func getProgressWheel(i int) (string, int) {

	wheels := []string{"-", "\\", "|", "/"}
	wheel_str := wheels[i]

	i = i + 1
	if i > 3 {
		i = 0
	}

	return wheel_str, i
}

func checkRunningStatus(progressWheel string) {
	if _, err := os.Stat("./ovm.pid"); !os.IsNotExist(err) {
		fmt.Printf("\033[2K\r")
		fmt.Printf("ovm is running %s", progressWheel)
	} else {
		fmt.Printf("\033[2K\r")
		fmt.Printf("ovm stopped normally")
	}
}

func main() {
	fmt.Println("starting daemon")

	i := 0
	for {
		time.Sleep(time.Second)
		progressWheel, newI := getProgressWheel(i)
		checkRunningStatus(progressWheel)
		i = newI
	}
}
