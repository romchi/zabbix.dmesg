package dmesg

import "fmt"

type errorString struct {
	str string
}

func (err *errorString) Error() string {
	return err.str
}

func levelValidate(level string) error {
	var availableLevels []string
	for item := range errLevels {
		availableLevels = append(availableLevels, item)
		if level == item {
			return nil
		}
	}
	return &errorString{fmt.Sprintf("Error: error level '%s' mismatch, available %s", level, availableLevels)}
}
