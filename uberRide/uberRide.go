package main

import (
	"fmt"
	"sync"
)

const (
	MaxDemocrats   = 16
	MaxRepublicans = 16
)

var wg sync.WaitGroup

type coordinator struct {
	republicans  int
	democrats    int
	demsWaiting  chan struct{}
	repubWaiting chan struct{}

	sync.Mutex
}

func NewCoordinator() *coordinator {
	return &coordinator{
		republicans:  0,
		democrats:    0,
		demsWaiting:  make(chan struct{}),
		repubWaiting: make(chan struct{}),
	}
}

func (c *coordinator) seatRepublican(id int) {
	rideLeader := false
	allrep := false
	c.Lock()
	c.republicans++
	if c.republicans == 4 {
		// Seat all republicans
		// release other 3 republicans
		<-c.repubWaiting
		<-c.repubWaiting
		<-c.repubWaiting

		// decrease count
		c.republicans -= 4

		// set rider leader to drive
		rideLeader = true
		allrep = true
	} else if c.republicans == 2 && c.democrats >= 2 {
		// Seat 2 rep 2 demo

		// release 2 demo 1 rep
		<-c.repubWaiting
		<-c.demsWaiting
		<-c.demsWaiting

		// decrease count
		c.republicans -= 2
		c.democrats -= 2

		// set rider leader to drive
		rideLeader = true
	} else {
		//fmt.Printf("Republican  id %d , demCounter = %d\n",id,c.republicans)

		c.Unlock()

		// If no match we wait in line (aquire semaphore)
		c.repubWaiting <- struct{}{}

	}

	// Print that we are seated
	seat(id)

	if rideLeader {
		drive(id, allrep)
		c.Unlock()
	}

}
func (c *coordinator) seatDemocrat(id int) {
	rideLeader := false
	allDem := false
	c.Lock()
	c.democrats++
	if c.democrats == 4 {
		// Seat all demo
		// release other 3 demo
		<-c.demsWaiting
		<-c.demsWaiting
		<-c.demsWaiting

		// decrease count
		c.democrats -= 4

		// set rider leader to drive
		rideLeader = true
		allDem = true
	} else if c.democrats == 2 && c.republicans >= 2 {
		// Seat 2 rep 2 demo

		// release 2 demo 1 rep
		<-c.demsWaiting
		<-c.repubWaiting
		<-c.repubWaiting

		// decrease count
		c.republicans -= 2
		c.democrats -= 2

		// set rider leader to drive
		rideLeader = true
	} else {
		//fmt.Printf("Democrat id %d , demCounter = %d\n",id,c.democrats)

		c.Unlock()

		// If no match we wait in line
		c.demsWaiting <- struct{}{}

	}

	// Print that we are seated
	seat(id)

	if rideLeader {
		drive(id, allDem)
		c.Unlock()
	}

}

func seat(id int) {
	fmt.Printf("Person with id %d seated\n", id)
}

func drive(id int, allsame bool) {
	fmt.Printf("Person with id %d drive, sametype = %t\n", id, allsame)
}

func main() {

	co := NewCoordinator()

	for r := 0; r < MaxRepublicans; r++ {
		wg.Add(1)
		go func(c *coordinator, r int) {
			defer wg.Done()
			c.seatRepublican(r)
		}(co, r)
	}

	for d := 0; d < MaxDemocrats; d++ {
		wg.Add(1)
		go func(c *coordinator, d int) {
			defer wg.Done()
			c.seatDemocrat(d)
		}(co, d)
	}
	wg.Wait()

}
