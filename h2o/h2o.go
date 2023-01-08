package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	MaxRoundOfGoroutines = 20
	WaitTime             = time.Millisecond
)

var wg sync.WaitGroup

func hydrogen(hchan chan chan struct{}) {
	defer wg.Done()
	time.Sleep(WaitTime)
	// Create channels so barrier can return response
	h := make(chan struct{})
	// Send this through another channel to the barrier
	hchan <- h
	// Received ACK from barrier
	<-h
	fmt.Print("H\n")
}

func oxygen(ochan chan chan struct{}) {
	defer wg.Done()
	time.Sleep(WaitTime)
	// Create channels so barrier can return response
	o := make(chan struct{})
	// Send this through another channel to the barrier
	ochan <- o
	// Received ACK from barrier
	<-o
	fmt.Print("O\n")
}

type barier struct {
	hchan chan chan struct{}
	ochan chan chan struct{}
}

func InitBarier(hchan, ochan chan chan struct{}) {
	for {
		h1ResponseChan := <-hchan
		h2ResponseChan := <-hchan
		o1ResponseChan := <-ochan

		h1ResponseChan <- struct{}{}
		h2ResponseChan <- struct{}{}
		o1ResponseChan <- struct{}{}

	}

}

func main() {

	hchan := make(chan chan struct{}, 2)
	ochan := make(chan chan struct{}, 1)

	go InitBarier(hchan, ochan)

	for i := 0; i < MaxRoundOfGoroutines; i++ {
		wg.Add(3)
		go hydrogen(hchan)
		go hydrogen(hchan)
		go oxygen(ochan)
	}
	wg.Wait()
}
