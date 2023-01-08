package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	NumberOfPhilosofers = 5
	CanEat              = 3
	EatingTime          = time.Millisecond
	ThinkingTime        = time.Millisecond
)

type stick struct {
	id int
	sync.Mutex
}

type philosopher struct {
	id                    int
	leftStick, rightStick *stick
}

func NewPhilosofer(id int, left, right *stick) *philosopher {
	return &philosopher{
		id:         id,
		leftStick:  left,
		rightStick: right,
	}
}

func (p *philosopher) dine() {
	for i := 0; i < CanEat; i++ {
		fmt.Printf("Philosofer %d is thinking\n", p.id)
		time.Sleep(ThinkingTime)
		// Get sticks
		p.leftStick.Lock()
		p.rightStick.Lock()

		// Eat for x time
		fmt.Printf("Philosofer %d is eating with fork %d %d\n", p.id, p.leftStick.id, p.rightStick.id)
		time.Sleep(EatingTime)
		fmt.Printf("Philosofer %d finished eating\n", p.id)

		// Return sticks
		p.leftStick.Unlock()
		p.rightStick.Unlock()

	}
}

func main() {

	var wg sync.WaitGroup

	chops := make([]*stick, NumberOfPhilosofers)

	for s := 0; s < NumberOfPhilosofers; s++ {
		chops[s] = &stick{id: s}
	}

	for i := 0; i < NumberOfPhilosofers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p := NewPhilosofer(id, chops[id], chops[(id+1)%NumberOfPhilosofers])
			p.dine()
		}(i)

	}
	wg.Wait()

}
