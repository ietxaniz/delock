// Package delock provides a deadlock detection mechanism for Go applications.
// It enhances the standard sync.Mutex with deadlock detection capabilities,
// including stack trace logging and timeout features.
package delock

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// Mutex extends the standard sync.Mutex with additional deadlock detection features.
// It includes mechanisms to set a timeout for lock acquisition and to capture stack traces
// when deadlocks are detected.
type Mutex struct {
	innerMu   sync.Mutex
	outerMu   sync.Mutex
	stackInfo map[int]stackInfoItem
	lastID    int
	Timeout   time.Duration
}

// SetTimeout sets the timeout duration for the Mutex. If a lock is not acquired within
// this duration, it is considered a deadlock situation.
// The timeout can be set globally for all Mutex instances via the DELOCK_TIMEOUT
// environment variable or 1000 milliseconds if the variable is not set.
func (m *Mutex) SetTimeout(timeout time.Duration) {
	m.innerMu.Lock()
	defer m.innerMu.Unlock()
	m.Timeout = timeout
}

// getTimeout retrieves the current timeout setting for the Mutex.
// If not set explicitly, it defaults to the value specified by the DELOCK_TIMEOUT
// environment variable or 1000 milliseconds if the variable is not set.
func (m *Mutex) getTimeout() time.Duration {
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

// Lock attempts to acquire the mutex lock. It returns a unique identifier for the lock
// and an error if a deadlock is detected (based on the timeout).
// The method captures a stack trace at the time of lock attempt to aid in debugging.
func (m *Mutex) Lock() (int, error) {
	timeout := m.getTimeout()
	stackData := getStackTrace()

	m.innerMu.Lock()
	if m.stackInfo == nil {
		m.stackInfo = make(map[int]stackInfoItem)
	}
	m.lastID++
	id := m.lastID
	m.stackInfo[id] = stackInfoItem{stackData: stackData, lock: WRITE_LOCK}
	m.innerMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	lockAcquired := make(chan struct{})
	go func() {
		m.outerMu.Lock()
		select {
		case lockAcquired <- struct{}{}:
		case <-ctx.Done():
			m.outerMu.Unlock()
		}
	}()

	select {
	case <-lockAcquired:
		return id, nil
	case <-ctx.Done():
		m.innerMu.Lock()
		delete(m.stackInfo, id)
		m.innerMu.Unlock()
		return -1, m.getErrorWithStackInfo()
	}
}

// Unlock releases the lock and removes its associated stack trace information.
func (m *Mutex) Unlock(infoID int) {
	m.innerMu.Lock()
	delete(m.stackInfo, infoID)
	m.innerMu.Unlock()

	m.outerMu.Unlock()
}

// getErrorWithStackInfo constructs an error with detailed stack trace information.
// This method is called when a deadlock is detected (i.e., lock acquisition times out).
func (m *Mutex) getErrorWithStackInfo() error {
	m.innerMu.Lock()
	defer m.innerMu.Unlock()
	report := createReport(m.stackInfo)
	report = fmt.Sprintf("\n\nDeadlock detected\n\n\n%s\n\n\n\n", report)
	return errors.New(report)
}
