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
	l.Execute()
}
