package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
)

var dataDir = os.Getenv("HOME") + "/.gchat"

func main() {
	config := &Config{}
	_, err := flags.Parse(config)
	if err != nil {
		os.Exit(1)
	}
	if err := config.Load(dataDir + "/chat.conf"); err != nil {
		log.Println(err)
	}

	chat := NewChat(dataDir, config)
	if err := chat.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
