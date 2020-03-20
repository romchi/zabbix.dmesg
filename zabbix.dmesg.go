package main

import (
	"fmt"
	"os"

	"github.com/romchi/zabbix.dmesg/dmesg"
)

func main() {
	messages, err := dmesg.NewLog()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	num, err := messages.CountLevel("warn")
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	fmt.Print(num)
}
