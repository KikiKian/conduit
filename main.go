package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := InitConduit(); err != nil {
		log.Fatalln("failed to init .conduit:", err)
	}

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println(dimStyle.Render("  usage: conduit <command>"))
		fmt.Println(dimStyle.Render("  commands: init, add, watch"))
		return
	}

	switch args[0] {
	case "init":
		runInit()

	case "add":
		runAdd()

	case "watch":
		var overrides []string
		for i := 1; i < len(args); i++ {
			if args[i] == "--file" && i+1 < len(args) {
				overrides = append(overrides, args[i+1])
				i++
			}
		}

		config, err := LoadConfig(overrides)
		if err != nil {
			log.Fatalln(errorStyle.Render("✗ " + err.Error()))
		}

		if err := Watch(config.Targets); err != nil {
			log.Fatalln(errorStyle.Render("✗ " + err.Error()))
		}

	default:
		fmt.Println(errorStyle.Render("  ✗ unknown command: " + args[0]))
		fmt.Println(dimStyle.Render("  commands: init, add, watch"))
	}
}
