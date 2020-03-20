package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/romchi/zabbix.dmesg/dmesg"
)

func main() {

	var countLevel string

	flag.StringVar(&countLevel, "l", "", "specify log level for count")
	flag.Usage = func() {
		fmt.Printf("Example usage of %s:\n", os.Args[0])
		fmt.Printf("    %s -l warn\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	getMessagesCount(countLevel)
}

func getMessagesCount(level string) {
	messages, err := dmesg.NewLog()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	num, err := messages.CountLevel(level)
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	fmt.Print(num)
}
