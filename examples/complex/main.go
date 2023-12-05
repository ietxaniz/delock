package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ietxaniz/delock"
)

type Child struct {
	mu     delock.RWMutex
	value  string
	id     int
	parent *Parent
}

type Parent struct {
	mu         delock.RWMutex
	items      []*Child
	lastID     int
	numChanges int
}

func (p *Parent) Add(value string) {
	lockID, err := p.mu.Lock()
	if err != nil {
		panic(err)
	}
	defer p.mu.Unlock(lockID)

	p.lastID++
	if p.items == nil {
		p.items = make([]*Child, 0)
	}
	p.items = append(p.items, &Child{value: value, id: p.lastID, parent: p})
}

func (p *Parent) GetChildByID(childID int) *Child {
	lockID, err := p.mu.RLock()
	if err != nil {
		panic(err)
	}
	defer p.mu.RUnlock(lockID)

	for _, item := range p.items {
		if item.GetID() == childID {
			return item
		}
	}

	return nil
}

func (p *Parent) GetNumChanges() int {
	lockID, err := p.mu.RLock()
	if err != nil {
		panic(err)
	}
	defer p.mu.RUnlock(lockID)

	return p.numChanges
}

func (p *Parent) IncreaseNumChanges() {
	lockID, err := p.mu.RLock()
	if err != nil {
		panic(err)
	}
	defer p.mu.RUnlock(lockID)

	p.numChanges++
}

func (c *Child) GetID() int {
	lockID, err := c.mu.RLock()
	if err != nil {
		panic(err)
	}
	defer c.mu.RUnlock(lockID)

	return c.id
}

func (c *Child) SetValue(value string) {
	lockID, err := c.mu.Lock()
	if err != nil {
		panic(err)
	}
	defer c.mu.Unlock(lockID)

	c.value = value
	c.parent.IncreaseNumChanges()
}

func main() {
	var p Parent
	if len(os.Args) < 3 {
		println("Usage: complex <n> <t>")
		return
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	t, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	go func() {
		for i := 0; i < n; i++ {
			time.Sleep(time.Duration(t) * time.Millisecond)
			p.Add(fmt.Sprintf("child%d", i))
			go func(id int, sleepMilliseconds int) {
				for {
					time.Sleep(time.Duration(sleepMilliseconds) * time.Millisecond)
					child := p.GetChildByID(id)
					if child == nil {
						log.Printf("Child %d not found", id)
					}
					if child != nil {
						child.SetValue(fmt.Sprintf("child%d", id))
					}
				}
			}(i+1, t)
		}
	}()

	lastReadedNumChanges := 0

	go func() {

		for {
			time.Sleep(1 * time.Second)
			lastReadedNumChanges = p.GetNumChanges()
		}
	}()

	for {
		time.Sleep(1 * time.Second)
		println(lastReadedNumChanges)
	}
}
