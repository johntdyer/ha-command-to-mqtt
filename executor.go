package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ExecuteCommandPeriodically runs a command at regular intervals
func ExecuteCommandPeriodically(cmd CommandConfig, clientID string) {
	frequency, err := time.ParseDuration(cmd.Frequency)
	if err != nil {
		logger.Errorf("Invalid frequency for command %s: %v", cmd.Name, err)
		frequency = 60 * time.Second // Default to 60 seconds
	}

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	// Execute immediately
	ExecuteCommand(cmd, clientID)

	// Then execute periodically
	for range ticker.C {
		ExecuteCommand(cmd, clientID)
	}
}

// ExecuteCommand executes a single command and publishes the result
func ExecuteCommand(cmd CommandConfig, clientID string) {
	logger.Debugf("Executing command: %s", cmd.Name)

	var result string

	// Default to local execution if target_host is not specified or is "local"
	if cmd.TargetHost != "" && cmd.TargetHost != "local" {
		// Execute command via SSH
		conn, exists := GetSSHConnection(cmd.TargetHost)
		if !exists {
			logger.Errorf("Target host %s not found for command %s", cmd.TargetHost, cmd.Name)
			result = fmt.Sprintf("ERROR: Target host %s not configured", cmd.TargetHost)
		} else {
			// Check if connection is still alive, reconnect if needed
			if !IsSSHConnectionAlive(conn) {
				logger.Warnf("SSH connection to %s is dead, reconnecting...", cmd.TargetHost)
				reconnectErr := ReconnectSSH(cmd.TargetHost)
				if reconnectErr != nil {
					logger.Errorf("Failed to reconnect to SSH host %s: %v", cmd.TargetHost, reconnectErr)
					result = fmt.Sprintf("ERROR: Failed to reconnect to SSH host: %v", reconnectErr)
				} else {
					// Get the new connection
					conn, _ = GetSSHConnection(cmd.TargetHost)
				}
			}

			if conn != nil {
				output, sshErr := ExecuteSSHCommand(conn, cmd.Command)
				if sshErr != nil {
					logger.Errorf("SSH command %s failed: %v", cmd.Name, sshErr)
					result = fmt.Sprintf("ERROR: %v", sshErr)
				} else {
					result = strings.TrimSpace(output)
				}
			}
		}
	} else {
		// Execute command locally
		result = executeLocalCommand(cmd)
	}

	// Publish result to MQTT
	PublishResult(cmd, result, clientID)
}

func executeLocalCommand(cmd CommandConfig) string {
	if strings.TrimSpace(cmd.Command) == "" {
		logger.Errorf("Empty command for %s", cmd.Name)
		return "ERROR: Empty command"
	}

	// Execute command using shell for proper interpretation of pipes, redirects, etc.
	execCmd := exec.Command("sh", "-c", cmd.Command)

	// Capture both stdout and stderr
	output, err := execCmd.CombinedOutput()

	if err != nil {
		outputStr := strings.TrimSpace(string(output))

		// Log the full output for debugging, especially if it's multi-line
		if strings.Contains(outputStr, "\n") {
			logger.Errorf("Command %s failed: %v\nFull output:\n%s", cmd.Name, err, outputStr)
		} else {
			logger.Errorf("Command %s failed: %v - (output: %s)", cmd.Name, err, outputStr)
		}

		// Return the actual output if available, otherwise return error
		if outputStr != "" {
			return fmt.Sprintf("ERROR: %s", outputStr)
		}
		return fmt.Sprintf("ERROR: %v", err)
	}

	result := strings.TrimSpace(string(output))

	// Log multi-line outputs for debugging
	if strings.Contains(result, "\n") {
		logger.Debugf("Command %s produced multi-line output:\n%s", cmd.Name, result)
	}

	return result
}