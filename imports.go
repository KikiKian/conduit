package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

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
			// extract variable name from comment e.g. "# conduit:import max_retries"
			parts := strings.Fields(trimmed)
			if len(parts) < 3 {
				outputLines = append(outputLines, line)
				continue
			}
			varName := parts[2]

			// load the value from the .conduit file
			entry, err := readEntry(varName)
			if err != nil {
				return fmt.Errorf("variable %s not found in .conduit: %w", varName, err)
			}

			// generate the python variable line
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
