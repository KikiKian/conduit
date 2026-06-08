package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type WatchTarget struct {
	Lang     string `json:"lang"`
	FilePath string `json:"filePath"`
}

type Config struct {
	Targets []WatchTarget `json:"targets"`
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

func Init() error {
	if _, err := os.Stat(".conduit"); os.IsNotExist(err) {
		return os.WriteFile(".conduit", []byte("[]"), 0644)
	}
	return nil
}

func AddEntry(name string, varType string, value any) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if _, ok := typeMap[varType]; !ok {
		return fmt.Errorf("unsupported type %q, supported: %s", varType, strings.Join(SupportedTypes(), ", "))
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

	entries = append(entries, Entry{Name: name, Type: varType, Value: value})

	out, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to serialize entries: %w", err)
	}

	return os.WriteFile(".conduit", out, 0644)
}

func GetEntry(name string) (Entry, error) {
	entries, err := ListEntries()
	if err != nil {
		return Entry{}, err
	}

	for _, e := range entries {
		if e.Name == name {
			return e, nil
		}
	}

	return Entry{}, fmt.Errorf("variable %q not found in .conduit", name)
}

func ListEntries() ([]Entry, error) {
	data, err := os.ReadFile(".conduit")
	if err != nil {
		return nil, fmt.Errorf("failed to read .conduit: %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse .conduit: %w", err)
	}

	return entries, nil
}

func CastValue(varType string, raw string) (any, error) {
	switch varType {
	case "int":
		n, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid int value %q", raw)
		}
		return n, nil
	case "bool":
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid bool value %q, use true or false", raw)
		}
		return b, nil
	case "string":
		return raw, nil
	default:
		return nil, fmt.Errorf("unsupported type %q", varType)
	}
}

func SupportedTypes() []string {
	keys := make([]string, 0, len(typeMap))
	for k := range typeMap {
		keys = append(keys, k)
	}
	return keys
}

func SupportedLangs() []string {
	return []string{"typescript", "python", "go", "java"}
}

func Save(targets []WatchTarget) error {
	if len(targets) == 0 {
		return fmt.Errorf("targets cannot be empty")
	}

	cfg := Config{Targets: targets}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	return os.WriteFile("conduit.config.json", data, 0644)
}

func Load(overrides []string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile("conduit.config.json")
	if err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	if err == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("failed to parse config: %w", err)
		}
	}

	for _, override := range overrides {
		parts := strings.SplitN(override, ":", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("invalid flag format %q, expected lang:filepath", override)
		}
		cfg.Targets = append(cfg.Targets, WatchTarget{
			Lang:     parts[0],
			FilePath: parts[1],
		})
	}

	if len(cfg.Targets) == 0 {
		return Config{}, fmt.Errorf("no targets found — add a conduit.config.json or use --file lang:filepath")
	}

	return cfg, nil
}

func GetTypeMap() map[string]map[string]string {
	return typeMap
}
