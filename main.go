package main // import "hello-channel"

import (
	"fmt"
	"time"
)

var triggerSubjob1 chan int
var finishSubjob1 chan int

var triggerSubjob2 chan int
var finishSubjob2 chan int

// subjob1 - subjob1
func subjob1() {
	fmt.Println("Waiting subjob1")

	for {
		if v, ok := <-triggerSubjob1; ok {
			// Receive - start subjob
			fmt.Println("Channel received", v)

			// Escape loop by -1 - finish subjob
			if v == -1 {
				break
			}

			for i := 0; i < 3; i++ {
				time.Sleep(180 * time.Millisecond)
				fmt.Println("Subjob", i)
			}
		}
	}

	// Notice main job that subjob is finished
	finishSubjob1 <- 0
}

// subjob2 - subjob2
func subjob2() {
	fmt.Println("Waiting subjob2")

JOBSDONE:
	for triggerSubjob2 != nil {
		select {
		case v := <-triggerSubjob2:
			// Receive - start subjob
			fmt.Println("Channel received", v)

			// Escape loop by -1 - finish subjob
			if v == -1 {
				finishSubjob2 <- 0
				break JOBSDONE
			}

			for i := 0; i < 3; i++ {
				time.Sleep(180 * time.Millisecond)
				fmt.Println("Subjob", i)
			}

		default:
		}
	}

	// Notice main job that subjob is finished
	finishSubjob1 <- 0
}

// main - main job
func main() {
	triggerSubjob1 = make(chan int)
	finishSubjob1 = make(chan int)
	triggerSubjob2 = make(chan int)
	finishSubjob2 = make(chan int)
	defer close(triggerSubjob1)
	defer close(finishSubjob1)
	defer close(triggerSubjob2)
	defer close(finishSubjob2)

	fmt.Println("<-Start")
	fmt.Println("")

	// Set subjob1 waiting
	fmt.Println("Ready to receive channel")
	go subjob1()

	// Set subjob2 waiting
	fmt.Println("Ready to receive channel")
	go subjob2()

	// main job
	fmt.Println("Main job start")
	for i := 0; i < 5; i++ {
		time.Sleep(180 * time.Millisecond)
		fmt.Println("Main job", i)
	}

	fmt.Println("----------------------------------------------------")

	// Trigger - start subjob
	for i := 0; i < 5; i++ {
		triggerSubjob1 <- i
	}
	// Send -1 to finish subjob
	triggerSubjob1 <- -1

	// Wait for subjob to finish
	<-finishSubjob1
	fmt.Println("Subjob1 finished")

	fmt.Println("----------------------------------------------------")

	// Trigger - start subjob
	for i := 0; i < 5; i++ {
		triggerSubjob2 <- i
	}
	// Send -1 to finish subjob
	triggerSubjob2 <- -1

	// Wait for subjob to finish
	<-finishSubjob2
	fmt.Println("Subjob2 finished")

	// All jobs done
	fmt.Println("")
	fmt.Println("<-End")
}
