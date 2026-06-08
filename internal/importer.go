package internal

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func ImportPython(filePath string) error {
	return processFile(filePath, "# conduit:import", func(varName string, entry Entry) string {
		pyType := typeMap[entry.Type]["python"]
		return fmt.Sprintf("%s: %s = %v", varName, pyType, entry.Value)
	})
}

func ImportTypescript(filePath string) error {
	return processFile(filePath, "// conduit:import", func(varName string, entry Entry) string {
		tsType := typeMap[entry.Type]["typescript"]
		return fmt.Sprintf("const %s: %s = %v", varName, tsType, entry.Value)
	})
}

func ImportGo(filePath string) error {
	return processFile(filePath, "// conduit:import", func(varName string, entry Entry) string {
		goType := typeMap[entry.Type]["go"]
		return fmt.Sprintf("var %s %s = %v", varName, goType, entry.Value)
	})
}

func processFile(filePath string, marker string, generate func(string, Entry) string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := bytes.Split(file, []byte("\n"))
	var outputLines [][]byte

	for _, line := range lines {
		trimmed := strings.TrimSpace(string(line))

		if strings.Contains(trimmed, marker) {
			parts := strings.Fields(trimmed)
			if len(parts) < 3 {
				outputLines = append(outputLines, line)
				continue
			}
			varName := parts[2]

			entry, err := GetEntry(varName)
			if err != nil {
				return fmt.Errorf("variable %s not found in .conduit: %w", varName, err)
			}

			outputLines = append(outputLines, []byte(generate(varName, entry)))
		} else {
			outputLines = append(outputLines, line)
		}
	}

	return os.WriteFile(filePath, bytes.Join(outputLines, []byte("\n")), 0644)
}
