package main

import (
	"log"
	"sync"
	"time"
)

type LockExample struct {
	mu    sync.Mutex
	Value int
}

func (l *LockExample) GetValue() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.Value
}
func (l *LockExample) Execute() {
	l.mu.Lock()
	defer l.mu.Unlock()

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
