package main

import (
	"fmt"
	"log"
	"os"

	"conduit/internal"
)

func main() {
	if err := internal.Init(); err != nil {
		log.Fatalln("failed to init .conduit:", err)
	}

	args := os.Args[1:]

	if len(args) == 0 {
		internal.PrintHelp()
		return
	}

	switch args[0] {
	case "init":
		internal.RunInit()

	case "add":
		internal.RunAdd()

	case "watch":
		var overrides []string
		for i := 1; i < len(args); i++ {
			if args[i] == "--file" && i+1 < len(args) {
				overrides = append(overrides, args[i+1])
				i++
			}
		}

		cfg, err := internal.Load(overrides)
		if err != nil {
			log.Fatalln(internal.ErrorStyle().Render("✗ " + err.Error()))
		}

		if err := internal.Watch(cfg.Targets); err != nil {
			log.Fatalln(internal.ErrorStyle().Render("✗ " + err.Error()))
		}

	default:
		fmt.Println(internal.ErrorStyle().Render("  ✗ unknown command: " + args[0]))
		fmt.Println(internal.DimStyle().Render("  commands: init, add, watch"))
	}
}
