package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/peterh/liner"
)

var defaultPorts = map[string]string{
	"mysql":  "3306",
	"oracle": "1521",
	"pgsql":  "5432",
}

var validTypes = map[string]bool{
	"mysql":  true,
	"oracle": true,
	"pgsql":  true,
}

// InteractiveSetup prompts the user for database connection info and returns a Config.
func InteractiveSetup() (*Config, error) {
	state := liner.NewLiner()
	defer state.Close()

	fmt.Println("Config file not found. Starting interactive setup:")

	// Database type
	dbType, err := promptChoice(state, "Database type [mysql/oracle/pgsql]", "", func(v string) bool {
		return validTypes[v]
	})
	if err != nil {
		return nil, err
	}

	// Host
	host, err := promptWithDefault(state, "Host", "127.0.0.1")
	if err != nil {
		return nil, err
	}

	// Port
	port, err := promptWithDefault(state, "Port", defaultPorts[dbType])
	if err != nil {
		return nil, err
	}
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("port must be a number: %s", port)
	}

	// Username
	username, err := promptRequired(state, "Username")
	if err != nil {
		return nil, err
	}

	// Password (masked)
	password, err := promptPassword(state, "Password")
	if err != nil {
		return nil, err
	}

	// Database name / service name
	dbNameLabel := "Database name"
	if dbType == "oracle" {
		dbNameLabel = "Service name"
	}
	dbName, err := promptRequired(state, dbNameLabel)
	if err != nil {
		return nil, err
	}

	// Assemble DSN
	dsn := buildDSN(dbType, username, password, host, port, dbName)

	cfg := &Config{
		Database: DatabaseConfig{
			Type:     dbType,
			DSN:      dsn,
			Username: username,
			Password: password,
		},
		Output: OutputConfig{
			Directory:     "./output",
			Format:        "csv",
			ShowInConsole: true,
			SaveToFile:    true,
		},
	}

	return cfg, nil
}

func buildDSN(dbType, username, password, host, port, dbName string) string {
	switch dbType {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)
	case "oracle":
		return fmt.Sprintf("oracle://%s:%s@%s:%s/%s", username, password, host, port, dbName)
	case "pgsql":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbName)
	default:
		return ""
	}
}

func promptWithDefault(state *liner.State, label, defaultVal string) (string, error) {
	prompt := fmt.Sprintf("%s [%s]: ", label, defaultVal)
	val, err := state.Prompt(prompt)
	if err != nil {
		return "", err
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return defaultVal, nil
	}
	return val, nil
}

func promptRequired(state *liner.State, label string) (string, error) {
	for {
		val, err := state.Prompt(label + ": ")
		if err != nil {
			return "", err
		}
		val = strings.TrimSpace(val)
		if val != "" {
			return val, nil
		}
		fmt.Println("This field is required")
	}
}

func promptChoice(state *liner.State, label, defaultVal string, valid func(string) bool) (string, error) {
	for {
		val, err := state.Prompt(label + ": ")
		if err != nil {
			return "", err
		}
		val = strings.TrimSpace(strings.ToLower(val))
		if val == "" && defaultVal != "" {
			return defaultVal, nil
		}
		if valid(val) {
			return val, nil
		}
		fmt.Println("Invalid input, please try again")
	}
}

func promptPassword(state *liner.State, label string) (string, error) {
	val, err := state.PasswordPrompt(label + ": ")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}
