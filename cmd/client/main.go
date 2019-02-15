package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func worker() {
	for {
		fmt.Println("Hello world!")
		time.Sleep(300 * time.Millisecond)
	}
}

func worker2() {
	for {
		fmt.Println("Foobar")
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {

	go worker()
	go worker2()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Finished, got signal:", s)
}
