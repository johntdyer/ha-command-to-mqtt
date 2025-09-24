package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize logger
	InitLogger()

	// Parse command line flags
	cliConfig := ParseFlags()

	// Set log level and format
	SetLogLevel(cliConfig.LogLevel, cliConfig.LogFormat)

	logger.Info("Starting HA Command to MQTT")

	// Load configuration
	config, err := LoadConfig(cliConfig.ConfigFile)
	if err != nil {
		logger.Fatal("Failed to load configuration:", err)
	}

	// Initialize SSH connections
	if err := InitSSHConnections(config); err != nil {
		logger.Fatal("Failed to initialize SSH connections:", err)
	}
	defer CloseSSHConnections()

	// Connect to MQTT
	if err := InitMQTT(&config.MQTT); err != nil {
		logger.Fatal("Failed to connect to MQTT:", err)
	}
	defer DisconnectMQTT()

	// Send discovery messages and start command execution
	for _, cmd := range config.Commands {
		SendDiscoveryMessage(cmd, config.MQTT.ClientID)
		go ExecuteCommandPeriodically(cmd, config.MQTT.ClientID)
	}

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Info("Shutting down...")
}