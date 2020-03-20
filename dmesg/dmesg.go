package dmesg

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"syscall"
	"time"
)

const dmesgBuferSize = 10
const dmesgReadAll = 3

var (
	errLevels = map[string]string{
		"emerg":  "0",
		"alert":  "1",
		"crit":   "2",
		"err":    "3",
		"warn":   "4",
		"notice": "5",
		"info":   "6",
		"debug":  "7",
	}

	// part 1: (?:<(\d+)>) 		// integer in <>
	// part 2: (?:\[\s1*?(\d+)   // first integer before . in [ ]
	// part 3: (\d+)\])			// integer in [ ] after .
	// part 4: (.*)\n?			// all string after [ ]
	bracketRegex = regexp.MustCompile(`(?:<(\d+)>)?(?:\[\s*?(\d+).(\d+)\])?(.*)\n?`)
)

// BootTime structure with system booting time
type BootTime struct {
	systemBootTime time.Time
}

// NewLog c
func NewLog() (*BootTime, error) {
	var err error
	log := BootTime{}
	info := syscall.Sysinfo_t{}
	bootTime := time.Now().Unix() - info.Uptime
	if err != nil {
		return nil, err
	}
	log.systemBootTime = time.Unix(bootTime, 0)
	return &log, nil
}

func (log *BootTime) readCurrent() ([]byte, error) {
	bufferSize, err := syscall.Klogctl(dmesgBuferSize, nil)
	if err != nil {
		return nil, fmt.Errorf("Error set size for log buffer: %s", err)
	}

	result := make([]byte, bufferSize)

	data, err := syscall.Klogctl(dmesgReadAll, result)
	if err != nil {
		return nil, fmt.Errorf("Error read messages: %s", err)
	}

	return result[:data], nil
}

// CountLevel read messages and return count for selected error level
func (log *BootTime) CountLevel(level string) (int, error) {
	err := levelValidate(level)
	if err != nil {
		return 0, err
	}

	buffer, err := log.readCurrent()
	if err != nil {
		return 0, err
	}

	buf := bytes.NewBuffer(buffer)
	var result int

	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Print(err)
			break
		}

		parts := bracketRegex.FindStringSubmatch(line)
		if parts != nil {
			if parts[1] == errLevels[level] {
				result++
			}
		}
	}
	return result, nil
}
