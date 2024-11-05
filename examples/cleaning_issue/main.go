package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/ietxaniz/delock"
)

type Item struct {
	mu    delock.RWMutex
	value int
}

func (i *Item) IncreaseValue() {
	lockID, err := i.mu.Lock()
	if err != nil {
		panic(err)
	}
	defer i.mu.Unlock(lockID)
	i.value = i.value + 1
}

func main() {
	var item Item
	item.IncreaseValue()
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			fmt.Printf("%d goroutines\n", runtime.NumGoroutine())
		}
	}()
	for {
		time.Sleep(100 * time.Millisecond)
		go func() {
			sleepMillis := 200 * rand.Float32()
			time.Sleep(time.Duration(sleepMillis) * time.Millisecond)
			item.IncreaseValue()
		}()
	}
}
