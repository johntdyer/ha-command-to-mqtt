package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// CLIConfig holds command line configuration
type CLIConfig struct {
	ConfigFile string
	LogLevel   string
	LogFormat  string
}

// ParseFlags parses command line flags and returns configuration
func ParseFlags() *CLIConfig {
	// Set default values
	viper.SetDefault("config", "config.yaml")
	viper.SetDefault("log-level", "info")
	viper.SetDefault("log-format", "text")

	// Set environment variable prefix
	viper.SetEnvPrefix("HA_MQTT")
	viper.AutomaticEnv()

	config := &CLIConfig{}

	// Parse command line flags manually (simple implementation)
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Handle help flags
		if arg == "-h" || arg == "--help" {
			printUsage()
			os.Exit(0)
		}

		// Handle version flag
		if arg == "-v" || arg == "--version" {
			fmt.Println("HA Command to MQTT v1.0.0")
			os.Exit(0)
		}

		// Handle config flags
		if arg == "-c" || arg == "--config" {
			if i+1 < len(args) {
				config.ConfigFile = args[i+1]
				i++ // Skip next argument as it's the value
			} else {
				fmt.Fprintf(os.Stderr, "Flag %s requires a value\n", arg)
				os.Exit(1)
			}
		} else if strings.HasPrefix(arg, "--config=") {
			config.ConfigFile = strings.TrimPrefix(arg, "--config=")
		}

		// Handle log level flags
		if arg == "-l" || arg == "--log-level" {
			if i+1 < len(args) {
				config.LogLevel = args[i+1]
				i++ // Skip next argument as it's the value
			} else {
				fmt.Fprintf(os.Stderr, "Flag %s requires a value\n", arg)
				os.Exit(1)
			}
		} else if strings.HasPrefix(arg, "--log-level=") {
			config.LogLevel = strings.TrimPrefix(arg, "--log-level=")
		}

		// Handle log format flags
		if arg == "-f" || arg == "--log-format" {
			if i+1 < len(args) {
				config.LogFormat = args[i+1]
				i++ // Skip next argument as it's the value
			} else {
				fmt.Fprintf(os.Stderr, "Flag %s requires a value\n", arg)
				os.Exit(1)
			}
		} else if strings.HasPrefix(arg, "--log-format=") {
			config.LogFormat = strings.TrimPrefix(arg, "--log-format=")
		}
	}

	// Set defaults from viper if not provided via flags
	if config.ConfigFile == "" {
		config.ConfigFile = viper.GetString("config")
	}
	if config.LogLevel == "" {
		config.LogLevel = viper.GetString("log-level")
	}
	if config.LogFormat == "" {
		config.LogFormat = viper.GetString("log-format")
	}

	return config
}

func printUsage() {
	fmt.Println("HA Command to MQTT v1.0.0 - Execute commands and publish results to MQTT for Home Assistant")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Printf("  %s [OPTIONS]\n", os.Args[0])
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -c, --config FILE     Configuration file path (default: config.yaml)")
	fmt.Println("  -l, --log-level LEVEL Log level: panic, fatal, error, warn, info, debug, trace (default: info)")
	fmt.Println("  -f, --log-format FMT  Log format: text, json, logfmt (default: text)")
	fmt.Println("  -v, --version         Show version information")
	fmt.Println("  -h, --help            Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Printf("  %s --config /path/to/config.yaml --log-level debug --log-format json\n", os.Args[0])
	fmt.Printf("  %s -c myconfig.yaml -l warn -f logfmt\n", os.Args[0])
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  HA_MQTT_CONFIG        Configuration file path")
	fmt.Println("  HA_MQTT_LOG_LEVEL     Log level")
	fmt.Println("  HA_MQTT_LOG_FORMAT    Log format")
}