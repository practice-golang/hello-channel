package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func sayWait(s string, ch chan string) {
	result := 0
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
		result = i
		// ch <- result // 이렇게 하니까 지금은 deadlock.
	}
	ch <- ("Done " + s + ": " + strconv.Itoa(result))
}

func main() {
	// c := make(chan string) // chan , 뒤에 값 넣으면 buffered 채널, 안 넣으면 unberfered 채널
	c := make(chan string, 1)
	defer close(c)

	var waiter sync.WaitGroup

	sayWait("Begin", c)
	fmt.Println(<-c)

	waiter.Add(1)
	go func() {
		sayWait("안녕 세상", c)
		fmt.Println(<-c)
		defer waiter.Done()
	}()

	waiter.Add(1)
	go func(s string) {
		sayWait(s, c)
		fmt.Println(<-c)
		defer waiter.Done()
	}("Hello world")

	waiter.Wait()
}
