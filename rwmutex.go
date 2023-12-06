package delock

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// RWMutex extends sync.RWMutex with additional deadlock detection features.
// It supports both read and write locks, with individual stack trace capture for each.
type RWMutex struct {
	innerMu   sync.Mutex
	outerMu   sync.RWMutex
	stackInfo map[int]stackInfoItem
	lastID    int
	Timeout   time.Duration
}

// SetTimeout sets the maximum duration to wait for a lock before considering it a deadlock.
// If a lock is not acquired within this duration, it's considered a deadlock situation.
func (m *RWMutex) SetTimeout(timeout time.Duration) {
	m.innerMu.Lock()
	defer m.innerMu.Unlock()
	m.Timeout = timeout
}

// getTimeout retrieves the current timeout. If not set explicitly, it defaults to the value
// specified in the DELOCK_TIMEOUT environment variable or 1000 milliseconds if not set.
func (m *RWMutex) getTimeout() time.Duration {
	m.innerMu.Lock()
	defer m.innerMu.Unlock()
	if m.Timeout == 0 {
		durationMs, err := strconv.Atoi(os.Getenv("DELOCK_TIMEOUT"))
		if err != nil {
			durationMs = 1000
		}
		m.Timeout = time.Duration(durationMs) * time.Millisecond
	}
	return m.Timeout
}

// Lock attempts to acquire a write lock. It returns a unique identifier for the lock
// and an error if a deadlock is detected based on the timeout.
// The method captures a stack trace at the time of lock attempt to aid in debugging.
func (m *RWMutex) Lock() (int, error) {
	timeout := m.getTimeout()
	stackData := ""
	bufferSize := 4096
	for {
		buffer := make([]byte, bufferSize)
		n := runtime.Stack(buffer, false)
		if n < bufferSize {
			stackData = string(buffer[:n])
			break
		}
		bufferSize *= 2
	}

	m.innerMu.Lock()
	if m.stackInfo == nil {
		m.stackInfo = make(map[int]stackInfoItem)
	}
	m.lastID = m.lastID + 1
	id := m.lastID
	m.stackInfo[id] = stackInfoItem{stackData: stackData, lock: WRITE_LOCK}
	m.innerMu.Unlock()

	timeoutChan := make(chan int)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- 0
	}()

	lockChan := make(chan int)
	go func() {
		m.outerMu.Lock()
		// Lock is acquired, but we need to check that timeout has not been fired before we return.
		lockNotifyChan := make(chan int)
		go func() {
			lockChan <- 0
		}()
		lockNotifyTimeoutChan := make(chan int)
		go func() {
			time.Sleep(100 * time.Millisecond)
		}()
		select {
		case <-lockNotifyChan:
			return
		case <-lockNotifyTimeoutChan:
			// Although lock has been acquired it is late so we release the acquired lock and finish the goroutine.
			m.outerMu.Unlock()
			return
		}
	}()

	select {
	case <-lockChan:
		return id, nil
	case <-timeoutChan:
		return -1, m.getErrorWithStackInfo()
	}
}

// Unlock releases the write lock and removes its associated stack trace information.
func (m *RWMutex) Unlock(infoID int) {
	m.innerMu.Lock()
	delete(m.stackInfo, infoID)
	m.innerMu.Unlock()

	m.outerMu.Unlock()
}

// getErrorWithStackInfo constructs an error with detailed stack trace information
// for all current locks. This method is called when a deadlock is detected.
func (m *RWMutex) getErrorWithStackInfo() error {
	m.innerMu.Lock()
	defer m.innerMu.Unlock()
	report := createReport(m.stackInfo)
	report = fmt.Sprintf("\n\nDeadlock detected\n\n\n%s\n\n\n\n", report)
	return errors.New(report)
}

// RLock attempts to acquire a read lock. Similar to Lock, it returns a unique identifier
// and an error if a deadlock is detected. It also captures the stack trace for debugging.
func (m *RWMutex) RLock() (int, error) {
	timeout := m.getTimeout()
	stackData := ""
	bufferSize := 4096
	for {
		buffer := make([]byte, bufferSize)
		n := runtime.Stack(buffer, false)
		if n < bufferSize {
			stackData = string(buffer[:n])
			break
		}
		bufferSize *= 2
	}

	m.innerMu.Lock()
	if m.stackInfo == nil {
		m.stackInfo = make(map[int]stackInfoItem)
	}
	m.lastID = m.lastID + 1
	id := m.lastID
	m.stackInfo[id] = stackInfoItem{stackData: stackData, lock: READ_LOCK}
	m.innerMu.Unlock()

	timeoutChan := make(chan int)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- 0
	}()

	lockChan := make(chan int)
	go func() {
		m.outerMu.RLock()
		// Lock is acquired, but we need to check that timeout has not been fired before we return.
		lockNotifyChan := make(chan int)
		go func() {
			lockChan <- 0
		}()
		lockNotifyTimeoutChan := make(chan int)
		go func() {
			time.Sleep(100 * time.Millisecond)
		}()
		select {
		case <-lockNotifyChan:
			return
		case <-lockNotifyTimeoutChan:
			// Although lock has been acquired it is late so we release the acquired lock and finish the goroutine.
			m.outerMu.RUnlock()
			return
		}
	}()

	select {
	case <-lockChan:
		return id, nil
	case <-timeoutChan:
		return -1, m.getErrorWithStackInfo()
	}
}

// RUnlock releases the read lock and removes its associated stack trace information.
func (m *RWMutex) RUnlock(infoID int) {
	m.innerMu.Lock()
	delete(m.stackInfo, infoID)
	m.innerMu.Unlock()

	m.outerMu.RUnlock()
}
