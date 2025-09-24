package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration structure
type Config struct {
	MQTT     MQTTConfig      `yaml:"mqtt"`
	SSH      SSHConfig       `yaml:"ssh,omitempty"`
	Commands []CommandConfig `yaml:"commands"`
}

// MQTTConfig holds MQTT broker configuration
type MQTTConfig struct {
	Broker   string `yaml:"broker"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ClientID string `yaml:"client_id"`
}

// SSHConfig holds SSH configuration
type SSHConfig struct {
	Hosts []SSHHost `yaml:"hosts,omitempty"`
}

// SSHHost represents an SSH host configuration
type SSHHost struct {
	Name       string `yaml:"name"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	KeyPath    string `yaml:"key_path,omitempty"`
	Password   string `yaml:"password,omitempty"`
	Timeout    string `yaml:"timeout,omitempty"`
}

// CommandConfig represents a command to be executed
type CommandConfig struct {
	Name           string `yaml:"name"`
	Command        string `yaml:"command"`
	Frequency      string `yaml:"frequency"` // duration string like "30s", "5m", "1h"
	DeviceClass    string `yaml:"device_class,omitempty"`
	Unit           string `yaml:"unit,omitempty"`
	Icon           string `yaml:"icon,omitempty"`
	TargetHost     string `yaml:"target_host,omitempty"` // Name of target host to execute command on ("local" or SSH host name)
	ForceUpdate    bool   `yaml:"force_update,omitempty"`
	StateClass     string `yaml:"state_class,omitempty"`
	EntityCategory string `yaml:"entity_category,omitempty"`
	ExpireAfter    int    `yaml:"expire_after,omitempty"`
}

// HomeAssistantDiscovery represents the HA discovery payload
type HomeAssistantDiscovery struct {
	Name              string `json:"name"`
	StateTopic        string `json:"state_topic"`
	UniqueID          string `json:"unique_id"`
	DeviceClass       string `json:"device_class,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	Icon              string `json:"icon,omitempty"`
	Device            Device `json:"device"`
	ForceUpdate       bool   `json:"force_update,omitempty"`
	StateClass        string `json:"state_class,omitempty"`
	EntityCategory    string `json:"entity_category,omitempty"`
	ExpireAfter       int    `json:"expire_after,omitempty"`
}

// Device represents the device information for HA
type Device struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Model        string   `json:"model"`
	Manufacturer string   `json:"manufacturer"`
}

// LoadConfig loads configuration from file or environment variables
func LoadConfig(configFile string) (*Config, error) {
	var config Config

	// First try to load from specified config file
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			logger.Infof("Loading configuration from: %s", configFile)
			return loadConfigFromYAML(configFile, &config)
		} else {
			logger.Infof("Config file %s not found", configFile)
		}
	}

	// Try default config.yaml if no config file specified or file not found
	if configFile == "config.yaml" || configFile == "" {
		if _, err := os.Stat("config.yaml"); err == nil {
			logger.Info("Loading configuration from: config.yaml")
			return loadConfigFromYAML("config.yaml", &config)
		}
	}

	// Fall back to environment variables
	logger.Info("Loading configuration from environment variables")
	return loadConfigFromEnv(&config)
}

func loadConfigFromYAML(filename string, config *Config) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", filename, err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %v", filename, err)
	}

	return config, nil
}

func loadConfigFromEnv(config *Config) (*Config, error) {
	// MQTT configuration from environment
	config.MQTT = MQTTConfig{
		Broker:   getEnvOrDefault("MQTT_BROKER", "localhost"),
		Port:     getEnvIntOrDefault("MQTT_PORT", 1883),
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
		ClientID: getEnvOrDefault("MQTT_CLIENT_ID", "ha-command-to-mqtt"),
	}

	// Commands from environment variables
	// Format: COMMAND_<NAME>=<command>
	// Optional: COMMAND_<NAME>_FREQUENCY=<duration>
	// Optional: COMMAND_<NAME>_DEVICE_CLASS=<class>
	// Optional: COMMAND_<NAME>_UNIT=<unit>
	// Optional: COMMAND_<NAME>_ICON=<icon>

	commands := make(map[string]CommandConfig)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		if strings.HasPrefix(key, "COMMAND_") {
			parsedKey := strings.TrimPrefix(key, "COMMAND_")

			if strings.Contains(parsedKey, "_FREQUENCY") {
				name := strings.TrimSuffix(parsedKey, "_FREQUENCY")
				cmd := commands[name]
				cmd.Name = name
				cmd.Frequency = value
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_DEVICE_CLASS") {
				name := strings.TrimSuffix(parsedKey, "_DEVICE_CLASS")
				cmd := commands[name]
				cmd.Name = name
				cmd.DeviceClass = value
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_UNIT") {
				name := strings.TrimSuffix(parsedKey, "_UNIT")
				cmd := commands[name]
				cmd.Name = name
				cmd.Unit = value
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_ICON") {
				name := strings.TrimSuffix(parsedKey, "_ICON")
				cmd := commands[name]
				cmd.Name = name
				cmd.Icon = value
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_TARGET_HOST") {
				name := strings.TrimSuffix(parsedKey, "_TARGET_HOST")
				cmd := commands[name]
				cmd.Name = name
				cmd.TargetHost = value
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_FORCE_UPDATE") {
				name := strings.TrimSuffix(parsedKey, "_FORCE_UPDATE")
				cmd := commands[name]
				cmd.Name = name
				cmd.ForceUpdate = strings.ToLower(value) == "true"
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_STATE_CLASS") {
				name := strings.TrimSuffix(parsedKey, "_STATE_CLASS")
				cmd := commands[name]
				cmd.Name = name
				cmd.StateClass = value
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_ENTITY_CATEGORY") {
				name := strings.TrimSuffix(parsedKey, "_ENTITY_CATEGORY")
				cmd := commands[name]
				cmd.Name = name
				cmd.EntityCategory = value
				commands[name] = cmd
			} else if strings.Contains(parsedKey, "_EXPIRE_AFTER") {
				name := strings.TrimSuffix(parsedKey, "_EXPIRE_AFTER")
				cmd := commands[name]
				cmd.Name = name
				if intValue, err := strconv.Atoi(value); err == nil {
					cmd.ExpireAfter = intValue
				}
				commands[name] = cmd
			} else {
				// This is the command itself
				cmd := commands[parsedKey]
				cmd.Name = parsedKey
				cmd.Command = value
				if cmd.Frequency == "" {
					cmd.Frequency = "60s" // Default frequency
				}
				commands[parsedKey] = cmd
			}
		}
	}

	// Convert map to slice
	for _, cmd := range commands {
		if cmd.Command != "" { // Only add commands that have an actual command
			config.Commands = append(config.Commands, cmd)
		}
	}

	if len(config.Commands) == 0 {
		return nil, fmt.Errorf("no commands configured")
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}