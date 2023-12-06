package delock

import (
	"regexp"
	"strings"
)

type parsedStackInfo struct {
	lockType        lockType
	goroutineNumber string
	status          string
	stackTrace      string
}

func normalizeStackTrace(stackTrace string) string {
	var normalizedLines []string
	lines := strings.Split(stackTrace, "\n")
	re := regexp.MustCompile(`:\d+`)
	found := false
	for _, line := range lines {
		if re.FindStringIndex(line) != nil {
			if !found {
				found = true
			} else {
				normalizedLines = append(normalizedLines, line)
			}
		}
	}
	return strings.Join(normalizedLines, "\n")
}

func parseStackData(stackData string, lockType lockType) parsedStackInfo {
	var goroutineNumber, status, stackTrace string

	re := regexp.MustCompile(`goroutine (\d+) \[(\w+)\]:\s*([\w\W]*)`)
	matches := re.FindStringSubmatch(stackData)

	if len(matches) >= 4 {
		goroutineNumber = matches[1]
		status = matches[2]
		stackTrace = normalizeStackTrace(matches[3])
	}

	return parsedStackInfo{
		goroutineNumber: goroutineNumber,
		status:          status,
		stackTrace:      stackTrace,
		lockType:        lockType,
	}
}
