package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"syscall"
	"time"
)

const dmesgBuferSize = 10
const dmesgReadAll = 3

var (
	errLevel = map[string]string{
		"emerg":  "0",
		"alert":  "1",
		"crit":   "2",
		"warn":   "4",
		"notice": "5",
		"info":   "6",
		"debug":  "7",
	}

	// part 1: (?:<(\d+)>) 		// integer in <>
	// part 2: (?:\[\s*?(\d+)   // first integer before . in [ ]
	// part 3: (\d+)\])			// integer in [ ] after .
	// part 4: (.*)\n?			// all string after [ ]
	bracketRegex = regexp.MustCompile(`(?:<(\d+)>)?(?:\[\s*?(\d+).(\d+)\])?(.*)\n?`)
)

// RawLog structure
type RawLog struct {
	systemBootTime time.Time
}

// Message encapsulates kernel ring buffer messages.  The timestamp is
// absolute, not relative to boot time.
type Message struct {
	Level     int64
	Timestamp time.Time
	Text      string
}

func main() {
	messages, err := NewLog()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	messages.byLevel("info")

	//d, err := messages.Messages()
	//if err != nil {
	//	fmt.Print(err)
	//	os.Exit(1)
	//}
	//for _, i := range d {
	//	fmt.Print(i)
	//}
}

// NewLog creates a new RawLog structure with boot time
func NewLog() (*RawLog, error) {
	var err error
	log := RawLog{}
	info := syscall.Sysinfo_t{}
	bootTime := time.Now().Unix() - info.Uptime
	if err != nil {
		return nil, err
	}
	log.systemBootTime = time.Unix(bootTime, 0)
	return &log, nil
}

// getCurrent retrieves the current contents of the kernel message ring buffer
func (log *RawLog) getCurrent() ([]byte, error) {
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

// Messages ss
func (log *RawLog) readMessages() ([]*Message, error) {
	buffer, err := log.getCurrent()
	if err != nil {
		return nil, err
	}
	messages, err := log.ProcessLog(buffer)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// Level ll
func (log *RawLog) byLevel(level string) {
	buffer, err := log.getCurrent()

	var result int

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	buf := bytes.NewBuffer(buffer)

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
			if parts[1] == errLevel[level] {
				result++
				//fmt.Print("level -> ", string(parts[1]))
				//fmt.Print("\n")
				//fmt.Print("time 1 -> ", string(parts[2]))
				//fmt.Print("\n")
				//fmt.Print("time 2 -> ", string(parts[3]))
				//fmt.Print("\n")
				//fmt.Print("message -> ", string(parts[4]))
				//fmt.Print("\n\n")
			}
		}
	}
	fmt.Print(result)
}

// ProcessLog processing log to struct
func (log *RawLog) ProcessLog(buffer []byte) ([]*Message, error) {
	//buf := bytes.NewBuffer(buffer)
	var result []*Message
	//for {
	//	line, err := buf.ReadString('\n')
	//	if err != nil {
	//		return result, err
	//	}
	//	//message, err := parseMessage(line)
	//	//if err != nil {
	//	//	return result, err
	//	//}
	//	//result = append(result, message)
	//}
	return result, nil
}

//func parseMessage(message string) (*Message, error) {
//	var level int64
//	var err error
//	var secs, nsecs int64
//
//	if message[1] == "" {
//		level = 6
//	} else {
//		level, err = strconv.ParseInt(message[1], 10, 64)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	if message[2] == "" {
//		secs = 0
//	} else {
//		secs, err = strconv.ParseInt(message[2], 10, 64)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	if message[3] == "" {
//		nsecs = 0
//	} else {
//		nsecs, err = strconv.ParseInt(message[3], 10, 64)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &Message{level, time.Unix(secs, nsecs), message[4]}, nil
//}

func dmesgRead() ([]byte, error) {
	bufferSize, err := syscall.Klogctl(dmesgBuferSize, nil)
	if err != nil {
		return nil, fmt.Errorf("Error query size log buffer: %s", err)
	}

	result := make([]byte, bufferSize)

	data, err := syscall.Klogctl(dmesgReadAll, result)
	if err != nil {
		return nil, fmt.Errorf("Error read messages: %s", err)
	}

	return result[:data], nil
}

//func parseDebugMessage(message string) (*Message, error) {
//parts := hmsRegex.FindStringSubmatch(message)
//if parts != nil {
//	return parseHMS(parts)
//}
//parts = bracketRegex.FindStringSubmatch(message)
//if parts != nil {
//	return parseBracketed(parts)
//}
//return nil, fmt.Errorf("Could not parse string %q", message)
//line, err := reader.ReadString('\n')
//if err != nil {
//	return nil, errors.New(fmt.Sprintf("Failed to parse line from reader [%s]", err))
//}
//}

func dmegsCount(level int) (int, error) {

	return 0, nil
}
