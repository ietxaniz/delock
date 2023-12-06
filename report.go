package delock

import (
	"fmt"
	"strings"
)

// lockType defines the type of lock - read or write.
type lockType int

// Constants representing the type of lock.
const (
	READ_LOCK  lockType = 1
	WRITE_LOCK lockType = 2
)

// stackInfoItem holds information about a lock, including its type and the stack data.
type stackInfoItem struct {
	lock      lockType
	stackData string
}

func createReport(items map[int]stackInfoItem) string {
	aggregatedData := make(map[string][]parsedStackInfo)

	for _, item := range items {
		parsedData := parseStackData(item.stackData, item.lock)
		key := fmt.Sprintf("%d_%s", item.lock, parsedData.stackTrace)
		aggregatedData[key] = append(aggregatedData[key], parsedData)
	}

	var report strings.Builder
	for _, data := range aggregatedData {
		// lockType := data[0].lock // Assuming all items in the group have the same lock type
		count := len(data)
		lockTypeString := "READ"
		if data[0].lockType == WRITE_LOCK {
			lockTypeString = "WRITE"
		}
		report.WriteString(fmt.Sprintf("\n\n* * * * * %s LOCK - %d occurrences * * * * * * * * \n%s",
			lockTypeString, count, data[0].stackTrace))
	}

	return report.String()
}
