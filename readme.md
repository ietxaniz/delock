# Delock

Delock is a Go library that enhances the standard `sync.Mutex` and `sync.RWMutex` with additional deadlock detection capabilities. It's designed to help developers identify and debug deadlock situations in Go applications by providing concise, grouped stack trace information and customizable timeout features.

## Features

- Extends `sync.Mutex` and `sync.RWMutex` with enhanced deadlock detection.
- Provides streamlined deadlock reports, grouping occurrences by lock type and focusing on key code lines for quicker analysis.
- Captures essential stack trace information when a deadlock is detected.
- Customizable timeout for lock acquisition.
- Minimal changes required for integration into existing Go code.

## Installation

To install Delock, use the following command:

```bash
go get github.com/ietxaniz/delock
```

## Usage

Replace your `sync.Mutex` or `sync.RWMutex` with `delock.Mutex` or `delock.RWMutex` respectively. Manage lock acquisition and release as shown in the examples below.

### Mutex Example

```go
package main

import (
    "log"
    "github.com/ietxaniz/delock"
)

func main() {
    mu := delock.Mutex{}

    id, err := mu.Lock()
    if err != nil {
        log.Fatal(err)
    }
    defer mu.Unlock(id)

    // Critical section
}
```

### RWMutex Example

```go
package main

import (
    "log"
    "github.com/ietxaniz/delock"
)

func main() {
    rwmu := delock.RWMutex{}

    go func() {
      id, err := rwmu.RLock()
      if err != nil {
        log.Printf("%s", err)
        return
      }
      defer rwmu.RUnlock(id)
      // Read critical section
    } ()

    id, err = rwmu.Lock()
    if err != nil {
        log.Fatal(err)
    }
    defer rwmu.Unlock(id)

    // Write critical section
}
```

## Configuration

You can set a global timeout for all `delock` mutexes by setting the `DELOCK_TIMEOUT` environment variable (in milliseconds). Alternatively, use the `SetTimeout` method to set the timeout duration programmatically.

## License

Delock is MIT licensed, as found in the [LICENSE](LICENSE) file.

## Contribution

Contributions are welcome! Feel free to submit pull requests, create issues, or propose new features.

## Further Reading

For more insights on how `delock` can be used to detect and debug deadlocks in Go applications, check out this article: [Go Deadlock Detection: Delock Library](https://dev.to/ietxaniz/go-deadlock-detection-delock-library-1eig).

## Release Notes

For the latest changes and updates, please refer to the [CHANGELOG](CHANGELOG.md).
