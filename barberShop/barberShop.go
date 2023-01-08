package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	HairCutDuration   = time.Millisecond
	SaloonMaxCapacity = 5
	MaxClients        = 20
)

var wg sync.WaitGroup

type barberShop struct {
	queue         []int
	ringBell      chan struct{}
	customersDone int
	mu            sync.Mutex
}

func NewBarberShop() *barberShop {
	return &barberShop{
		queue:         make([]int, 0),
		ringBell:      make(chan struct{}),
		customersDone: 0,
	}
}

func (bs *barberShop) isFull() bool {
	return len(bs.queue) == SaloonMaxCapacity
}

func (bs *barberShop) isEmpty() bool {
	return len(bs.queue) == 0
}

func (bs *barberShop) startDayForBarber() {
	for {
		bs.mu.Lock()
		if bs.customersDone == MaxClients {
			bs.mu.Unlock()
			wg.Done()
			return
		}
		if len(bs.queue) == 0 {
			fmt.Println("Barber started sleeping")
			if bs.customersDone == MaxClients {
				bs.mu.Unlock()
				wg.Done()
				return
			}
			bs.mu.Unlock()

			<-bs.ringBell
			fmt.Println("Barber woke up")
			continue
		}
		for len(bs.queue) != 0 {
			poppedVal := bs.queue[0]
			bs.queue = bs.queue[1:]
			fmt.Printf("Barber starting to cut customer with id %d\n", poppedVal)
			time.Sleep(HairCutDuration)
			fmt.Printf("Barber finished cutting customer with id %d\n", poppedVal)
		}
		bs.mu.Unlock()

	}
}

func (bs *barberShop) costumerEnters(id int) {
	shouldWakeBarber := false
	bs.mu.Lock()
	if bs.isFull() {
		fmt.Printf("Room full, customer %d leaves\n", id)
		bs.mu.Unlock()
		return
	}
	if bs.isEmpty() {
		shouldWakeBarber = true
	}
	fmt.Printf("Customer %d waiting\n", id)
	bs.queue = append(bs.queue, id)
	bs.mu.Unlock()
	if shouldWakeBarber {
		fmt.Printf("Bell ring by customer %d\n", id)
		bs.ringBell <- struct{}{}
	}

}

func main() {

	bs := NewBarberShop()
	wg.Add(1)

	go bs.startDayForBarber()

	for id := 0; id < MaxClients; id++ {
		wg.Add(1)
		go func(bs *barberShop, id int) {
			defer wg.Done()
			bs.costumerEnters(id)
		}(bs, id)
	}

	wg.Wait()

}
