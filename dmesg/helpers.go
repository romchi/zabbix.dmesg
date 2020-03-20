package dmesg

type errorString struct {
	str string
}

func (err *errorString) Error() string {
	return err.str
}

func levelValidate(level string) error {
	for item := range errLevels {
		if level == item {
			return nil
		}
	}
	return &errorString{"Error: error level mismatch"}
}
