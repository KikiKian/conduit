package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatchTarget struct {
	Lang     string `json:"lang"`
	FilePath string `json:"filePath"`
}

func Watch(targets []WatchTarget) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	err = watcher.Add(".conduit")
	if err != nil {
		return fmt.Errorf("failed to watch .conduit: %w", err)
	}

	if isTerminal() {
		printWatchHeader(targets)
	} else {
		fmt.Println("conduit watching for changes...")
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				ts := time.Now().Format("15:04:05")
				for _, target := range targets {
					if err := runImport(target); err != nil {
						if isTerminal() {
							printWatchError(ts, target.FilePath, err)
						} else {
							fmt.Printf("[%s] error updating %s: %v\n", ts, target.FilePath, err)
						}
					} else {
						if isTerminal() {
							printWatchEvent(ts, target.FilePath, target.Lang)
						} else {
							fmt.Printf("[%s] updated %s (%s)\n", ts, target.FilePath, target.Lang)
						}
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("watcher error: %v\n", err)
		}
	}
}

func runImport(target WatchTarget) error {
	switch target.Lang {
	case "python":
		return ImportPython(target.FilePath)
	case "typescript":
		return ImportTypescript(target.FilePath)
	case "go":
		return ImportGo(target.FilePath)
	default:
		return fmt.Errorf("unsupported language: %s", target.Lang)
	}
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
