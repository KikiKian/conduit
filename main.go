package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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
