package main

import (
	"fmt"
	"sync"
)

const (
	MaxInput = 30
)

type SecureVar struct {
	i    int
	wg   *sync.WaitGroup
	cond *sync.Cond
}

func NewSecureVar(wg *sync.WaitGroup) *SecureVar {
	mutex := new(sync.Mutex)
	cond := sync.NewCond(mutex)
	return &SecureVar{i: 1, cond: cond, wg: wg}
}

func (sv *SecureVar) buzz() {
	defer sv.wg.Done()
	for sv.i <= MaxInput {
		sv.cond.L.Lock()
		if sv.i%5 == 0 && sv.i%3 != 0 {
			fmt.Print("Buzz\n")
			sv.i++
			sv.cond.Broadcast()
			sv.cond.L.Unlock()
			continue
		}
		
		sv.cond.Wait()
		sv.cond.L.Unlock()

	}

}
func (sv *SecureVar) fizz() {
	defer sv.wg.Done()
	for sv.i <= MaxInput {
		sv.cond.L.Lock()

		if sv.i%5 != 0 && sv.i%3 == 0 {
			fmt.Print("Fizz\n")
			sv.i++
			sv.cond.Broadcast()
			sv.cond.L.Unlock()
			continue
		}
		
		sv.cond.Wait()
		sv.cond.L.Unlock()

	}

}
func (sv *SecureVar) fizzbuzz() {
	defer sv.wg.Done()
	for sv.i <= MaxInput {
		sv.cond.L.Lock()

		if sv.i%5 == 0 && sv.i%3 == 0 {			
			fmt.Print("FizzBuzz\n")
			sv.i++
			sv.cond.Broadcast()
			sv.cond.L.Unlock()
			continue
		}
		sv.cond.Wait()
		sv.cond.L.Unlock()

	}

}
func (sv *SecureVar) number() {
	defer sv.wg.Done()
	for sv.i <= MaxInput {
		sv.cond.L.Lock()
		if sv.i%5 != 0 && sv.i%3 != 0 {
			fmt.Printf("%d\n", sv.i)
			sv.i++
			sv.cond.Broadcast()
			sv.cond.L.Unlock()
			continue
		}
		
		sv.cond.Wait()
		sv.cond.L.Unlock()

	}
}

func main() {
	var wg sync.WaitGroup

	sv := NewSecureVar(&wg)

	wg.Add(4)
	go sv.fizz()
	go sv.buzz()
	go sv.fizzbuzz()
	go sv.number()
	wg.Wait()

}
