package main

import (
	"log"
	"time"

	"github.com/ietxaniz/delock"
)

type LockExample struct {
	mu    delock.Mutex
	Value int
}

func (l *LockExample) GetValue() int {
	lockID, err := l.mu.Lock()
	if err != nil {
		panic(err)
	}
	defer l.mu.Unlock(lockID)

	return l.Value
}
func (l *LockExample) Execute() {
	lockID, err := l.mu.Lock()
	if err != nil {
		panic(err)
	}
	defer l.mu.Unlock(lockID)

	time.Sleep(100 * time.Millisecond)
	l.Value++
	log.Printf("value: %d", l.Value)
	l.Execute()
}

func main() {
	var l LockExample
	go func() {
		l.Execute()
	}()
	readValueChan := make(chan int)
	go func() {
		readValueChan <- l.GetValue()
	}()
	for {
		timeoutChan := make(chan int)
		go func() {
			time.Sleep(1 * time.Second)
			timeoutChan <- 0
		}()

		select {
		case <-timeoutChan:
			log.Printf("waiting...")
		case <-readValueChan:
			log.Printf("value: %d", l.GetValue())
			return
		}
	}
}
