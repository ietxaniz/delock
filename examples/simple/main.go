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

func (l *LockExample) Execute() {
	id, err := l.mu.Lock()
	if err != nil {
		log.Fatal(err)
	}
	defer l.mu.Unlock(id)

	time.Sleep(100 * time.Millisecond)
	l.Value++
	log.Printf("value: %d", l.Value)
	l.Execute()
}

func main() {
	var l LockExample
	l.Execute()
}
