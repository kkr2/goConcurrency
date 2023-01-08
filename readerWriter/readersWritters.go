package readerwriter

import (
	"fmt"
	"sync"
)

const (
	NumberOfReaders = 200
	NumberOfWriters = 50
)

type rwCustomLock struct {
	readersCnt int
	noReaders  chan struct{}
	data       int
	mu         sync.Mutex
}

func NewCL() *rwCustomLock {
	return &rwCustomLock{
		readersCnt: 0,
		noReaders:  make(chan struct{}, 1),
		data:       0,
	}
}

func (rwl *rwCustomLock) writer(threadNr int) {
	//Aquire no readers semaphore
	rwl.noReaders <- struct{}{}

	// Write data (incerement data)
	rwl.data++
	fmt.Printf("Writer %d modified data to %d\n", threadNr, rwl.data)

	// Release no readers semaphore
	<-rwl.noReaders
}

func (rwl *rwCustomLock) reader(threadNr int) {
	// Reader acquire the lock before modifying count
	rwl.mu.Lock()
	rwl.readersCnt++
	// If this is the first reader we will aquire semaphore
	if rwl.readersCnt == 1 {
		rwl.noReaders <- struct{}{}
	}
	rwl.mu.Unlock()

	// Read section
	fmt.Printf("Reader %d read data %d\n", threadNr, rwl.data)

	// Reader acquire the lock before modifying count
	rwl.mu.Lock()
	rwl.readersCnt--
	// If this is the last writer we will release semaphore
	if rwl.readersCnt == 0 {
		<-rwl.noReaders
	}
	rwl.mu.Unlock()

}

func main() {
	rwLock := NewCL()
	var wg sync.WaitGroup

	for i := 0; i < NumberOfWriters; i++ {
		wg.Add(1)
		go func(goroutinNr int, rwLock *rwCustomLock) {
			defer wg.Done()
			rwLock.writer(goroutinNr)

		}(i, rwLock)

	}

	for i := 0; i < NumberOfReaders; i++ {
		wg.Add(1)
		go func(goroutinNr int, rwLock *rwCustomLock) {
			defer wg.Done()
			rwLock.reader(goroutinNr)

		}(i, rwLock)
	}

	wg.Wait()
}
