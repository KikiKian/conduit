package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func cmdInit() error {
	fmt.Println("welcome to conduit!")
	fmt.Println("let's set up your project.\n")

	scanner := bufio.NewScanner(os.Stdin)
	var targets []WatchTarget

	for {
		fmt.Print("add a target language (python, typescript, go) or press enter to finish: ")
		scanner.Scan()
		lang := strings.TrimSpace(scanner.Text())

		if lang == "" {
			break
		}

		if lang != "python" && lang != "typescript" && lang != "go" {
			fmt.Printf("unsupported language %q, supported: python, typescript, go\n", lang)
			continue
		}

		fmt.Printf("file path for %s: ", lang)
		scanner.Scan()
		filePath := strings.TrimSpace(scanner.Text())

		if filePath == "" {
			fmt.Println("file path cannot be empty")
			continue
		}

		targets = append(targets, WatchTarget{
			Lang:     lang,
			FilePath: filePath,
		})

		fmt.Printf("added %s → %s\n\n", lang, filePath)
	}

	if len(targets) == 0 {
		return fmt.Errorf("no targets added, exiting")
	}

	config := Config{Targets: targets}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	err = os.WriteFile("conduit.config.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Println("\nconduit.config.json created successfully!")
	fmt.Println("run conduit watch to start watching for changes.")

	return nil
}

func cmdAdd(args []string) error {
	var name, varType, varValue string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--name":
			if i+1 < len(args) {
				name = args[i+1]
				i++
			}
		case "--type":
			if i+1 < len(args) {
				varType = args[i+1]
				i++
			}
		case "--value":
			if i+1 < len(args) {
				varValue = args[i+1]
				i++
			}
		}
	}

	scanner := bufio.NewScanner(os.Stdin)

	if name == "" {
		fmt.Print("variable name: ")
		scanner.Scan()
		name = strings.TrimSpace(scanner.Text())
	}

	if varType == "" {
		fmt.Print("variable type (int, string, bool): ")
		scanner.Scan()
		varType = strings.TrimSpace(scanner.Text())
	}

	if varValue == "" {
		fmt.Print("variable value: ")
		scanner.Scan()
		varValue = strings.TrimSpace(scanner.Text())
	}

	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if _, ok := typeMap[varType]; !ok {
		return fmt.Errorf("unsupported type %q, supported: int, string, bool", varType)
	}

	var typedValue any
	switch varType {
	case "int":
		n, err := strconv.Atoi(varValue)
		if err != nil {
			return fmt.Errorf("invalid int value %q", varValue)
		}
		typedValue = n
	case "bool":
		b, err := strconv.ParseBool(varValue)
		if err != nil {
			return fmt.Errorf("invalid bool value %q, use true or false", varValue)
		}
		typedValue = b
	case "string":
		typedValue = varValue
	}

	var entries []Entry
	data, err := os.ReadFile(".conduit")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .conduit: %w", err)
	}

	if err == nil {
		if err := json.Unmarshal(data, &entries); err != nil {
			return fmt.Errorf("failed to parse .conduit: %w", err)
		}
	}

	for _, e := range entries {
		if e.Name == name {
			return fmt.Errorf("variable %q already exists in .conduit", name)
		}
	}

	entries = append(entries, Entry{
		Name:  name,
		Type:  varType,
		Value: typedValue,
	})

	out, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to serialize entries: %w", err)
	}

	if err := os.WriteFile(".conduit", out, 0644); err != nil {
		return fmt.Errorf("failed to write .conduit: %w", err)
	}

	fmt.Printf("\nadded %s (%s) = %v to .conduit\n", name, varType, typedValue)
	return nil
}
