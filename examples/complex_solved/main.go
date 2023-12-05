package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Child struct {
	mu     sync.RWMutex
	value  string
	id     int
	parent *Parent
}

type Parent struct {
	mu         sync.RWMutex
	items      []*Child
	lastID     int
	numChanges int
}

func (p *Parent) Add(value string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	time.Sleep(1 * time.Millisecond)
	p.lastID++
	if p.items == nil {
		p.items = make([]*Child, 0)
	}
	p.items = append(p.items, &Child{value: value, id: p.lastID, parent: p})
}

func (p *Parent) GetChildByID(childID int) *Child {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, item := range p.items {
		if item.GetID() == childID {
			return item
		}
	}

	return nil
}

func (p *Parent) GetNumChanges() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.numChanges
}

func (p *Parent) IncreaseNumChanges() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	time.Sleep(1 * time.Millisecond)

	p.numChanges++
}

func (c *Child) GetID() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	time.Sleep(1 * time.Millisecond)

	return c.id
}

func (c *Child) setValueSafe(value string) *Parent {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value = value
	return c.parent
}

func (c *Child) SetValue(value string) {
	parent := c.setValueSafe(value)
	parent.IncreaseNumChanges()
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
					if child != nil {
						child.SetValue(fmt.Sprintf("child%d", id))
					}
				}
			}(i, t)
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
