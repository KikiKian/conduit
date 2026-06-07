package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Entry struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

var typeMap = map[string]map[string]string{
	"int": {
		"typescript": "number",
		"python":     "int",
		"java":       "int",
		"go":         "int",
	},
	"string": {
		"typescript": "string",
		"python":     "str",
		"java":       "String",
		"go":         "string",
	},
	"bool": {
		"typescript": "boolean",
		"python":     "bool",
		"java":       "boolean",
		"go":         "bool",
	},
}

func conduitFile() {
	os.Create(".conduit")
}

func newEntry(name string, variable any) error {
	entry := Entry{
		Name:  name,
		Type:  reflect.TypeOf(variable).String(),
		Value: variable,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to serialize entry: %w", err)
	}

	err = os.WriteFile(name+".conduit", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func importPython(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := bytes.Split(file, []byte("\n"))
	var outputLines [][]byte

	for _, line := range lines {
		trimmed := strings.TrimSpace(string(line))

		if strings.Contains(trimmed, "# conduit:import") {

			parts := strings.Fields(trimmed)
			if len(parts) < 3 {
				outputLines = append(outputLines, line)
				continue
			}
			varName := parts[2]

			entry, err := readEntry(varName)
			if err != nil {
				return fmt.Errorf("variable %s not found in .conduit: %w", varName, err)
			}

			pyType := typeMap[entry.Type]["python"]
			generated := fmt.Sprintf("%s: %s = %v", varName, pyType, entry.Value)
			outputLines = append(outputLines, []byte(generated))
		} else {
			outputLines = append(outputLines, line)
		}
	}

	output := bytes.Join(outputLines, []byte("\n"))
	return os.WriteFile(filePath, output, 0644)
}

func readEntry(name string) (Entry, error) {
	data, err := os.ReadFile(".conduit")
	if err != nil {
		return Entry{}, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return Entry{}, err
	}

	for _, e := range entries {
		if e.Name == name {
			return e, nil
		}
	}

	return Entry{}, fmt.Errorf("entry not found")
}
