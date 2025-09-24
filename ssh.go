package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHConnection holds an active SSH connection
type SSHConnection struct {
	client *ssh.Client
	config SSHHost
}

var sshConnections map[string]*SSHConnection

// InitSSHConnections initializes SSH connections based on configuration
func InitSSHConnections(config *Config) error {
	sshConnections = make(map[string]*SSHConnection)

	// Check if SSH configuration exists and has hosts
	if len(config.SSH.Hosts) == 0 {
		logger.Info("No SSH hosts configured, only local commands will be available")
		return nil
	}

	logger.Infof("Initializing %d SSH connection(s)...", len(config.SSH.Hosts))

	successCount := 0
	for _, host := range config.SSH.Hosts {
		conn, err := createSSHConnection(host)
		if err != nil {
			logger.Errorf("Failed to connect to SSH host %s: %v", host.Name, err)
			continue // Don't fail completely, just log and continue
		}
		sshConnections[host.Name] = conn
		logger.Infof("Connected to SSH host: %s", host.Name)
		successCount++
	}

	logger.Infof("Successfully connected to %d/%d SSH hosts", successCount, len(config.SSH.Hosts))
	return nil
}

func createSSHConnection(hostConfig SSHHost) (*SSHConnection, error) {
	var authMethods []ssh.AuthMethod

	// Try key-based authentication first
	if hostConfig.KeyPath != "" {
		logger.Debugf("Loading specified SSH key for host %s: %s", hostConfig.Name, hostConfig.KeyPath)
		key, err := loadSSHKey(hostConfig.KeyPath)
		if err != nil {
			logger.Errorf("Failed to load specified SSH key for host %s from %s: %v", hostConfig.Name, hostConfig.KeyPath, err)
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(key))
			logger.Debugf("Successfully loaded specified SSH key for host %s", hostConfig.Name)
		}
	}

	// Add password authentication if provided
	if hostConfig.Password != "" {
		authMethods = append(authMethods, ssh.Password(hostConfig.Password))
	}

	// Try additional authentication methods if no explicit key/password configured
	if len(authMethods) == 0 {
		logger.Debugf("No explicit authentication configured for host %s, trying alternative methods...", hostConfig.Name)

		// First, try SSH agent
		logger.Debugf("Trying SSH agent for host %s", hostConfig.Name)
		if agentAuth := trySSHAgent(hostConfig.Name); agentAuth != nil {
			authMethods = append(authMethods, agentAuth)
			logger.Debugf("SSH agent authentication available for host %s", hostConfig.Name)
		}

		// If SSH agent not available, try default SSH keys
		if len(authMethods) == 0 {
			logger.Debugf("SSH agent not available, trying default SSH keys for host %s", hostConfig.Name)
			homeDir, err := os.UserHomeDir()
			if err != nil {
				logger.Warnf("Could not determine home directory for host %s: %v", hostConfig.Name, err)
			} else {
				defaultKeys := []string{
					filepath.Join(homeDir, ".ssh", "id_rsa"),
					filepath.Join(homeDir, ".ssh", "id_ed25519"),
					filepath.Join(homeDir, ".ssh", "id_ecdsa"),
				}

				keysFound := false
				for _, keyPath := range defaultKeys {
					logger.Debugf("Trying SSH key: %s", keyPath)
					if key, err := loadSSHKey(keyPath); err == nil {
						authMethods = append(authMethods, ssh.PublicKeys(key))
						logger.Debugf("Successfully loaded SSH key: %s", keyPath)
						keysFound = true
						break
					} else {
						logger.Debugf("Could not load SSH key %s: %v", keyPath, err)
					}
				}

				if !keysFound {
					logger.Warnf("No SSH keys found in default locations for host %s. Checked: %v", hostConfig.Name, defaultKeys)
				}
			}
		}
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication methods available for host %s (no SSH keys found, no SSH agent available, and no password provided)", hostConfig.Name)
	}

	// Set default timeout
	timeout := 30 * time.Second
	if hostConfig.Timeout != "" {
		if t, err := time.ParseDuration(hostConfig.Timeout); err == nil {
			timeout = t
		}
	}

	// Set default port
	port := 22
	if hostConfig.Port != 0 {
		port = hostConfig.Port
	}

	sshConfig := &ssh.ClientConfig{
		User:            hostConfig.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
		Timeout:         timeout,
	}

	addr := fmt.Sprintf("%s:%d", hostConfig.Host, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", addr, err)
	}

	return &SSHConnection{
		client: client,
		config: hostConfig,
	}, nil
}

func loadSSHKey(keyPath string) (ssh.Signer, error) {
	// Check if the key file exists first
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("SSH key file does not exist: %s", keyPath)
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key file %s: %v", keyPath, err)
	}

	// Try to parse as encrypted key first
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		// If that fails, it might be encrypted - for now we don't support passphrase
		return nil, fmt.Errorf("failed to parse SSH private key %s (encrypted keys with passphrases not supported): %v", keyPath, err)
	}

	return signer, nil
}

// CloseSSHConnections closes all SSH connections
func CloseSSHConnections() {
	for name, conn := range sshConnections {
		if conn != nil && conn.client != nil {
			conn.client.Close()
			logger.Infof("Closed SSH connection to: %s", name)
		}
	}
}

// ExecuteSSHCommand executes a command on an SSH connection
func ExecuteSSHCommand(conn *SSHConnection, command string) (string, error) {
	session, err := conn.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// IsSSHConnectionAlive tests if an SSH connection is still alive
func IsSSHConnectionAlive(conn *SSHConnection) bool {
	if conn == nil || conn.client == nil {
		return false
	}

	// Test the connection by sending a simple keepalive
	_, _, err := conn.client.SendRequest("keepalive@openssh.com", false, nil)
	return err == nil
}

// GetSSHConnection returns the SSH connection for a given host name
func GetSSHConnection(hostName string) (*SSHConnection, bool) {
	conn, exists := sshConnections[hostName]
	return conn, exists
}

// ReconnectSSH attempts to reconnect to an SSH host
func ReconnectSSH(hostName string) error {
	conn, exists := sshConnections[hostName]
	if !exists {
		return fmt.Errorf("SSH host %s not found", hostName)
	}

	newConn, err := createSSHConnection(conn.config)
	if err != nil {
		return err
	}

	sshConnections[hostName] = newConn
	return nil
}

// trySSHAgent attempts to connect to SSH agent and return authentication method
func trySSHAgent(hostName string) ssh.AuthMethod {
	// Try to connect to SSH agent via SSH_AUTH_SOCK environment variable
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		logger.Debugf("SSH_AUTH_SOCK not set for host %s, SSH agent not available", hostName)
		return nil
	}

	// Connect to SSH agent
	conn, err := net.Dial("unix", sshAuthSock)
	if err != nil {
		logger.Debugf("Failed to connect to SSH agent for host %s: %v", hostName, err)
		return nil
	}

	// Create agent client
	agentClient := agent.NewClient(conn)

	// Get list of keys from agent
	keys, err := agentClient.List()
	if err != nil {
		logger.Debugf("Failed to list keys from SSH agent for host %s: %v", hostName, err)
		conn.Close()
		return nil
	}

	if len(keys) == 0 {
		logger.Debugf("No keys available in SSH agent for host %s", hostName)
		conn.Close()
		return nil
	}

	logger.Debugf("Found %d key(s) in SSH agent for host %s", len(keys), hostName)

	// Return SSH agent authentication method
	return ssh.PublicKeysCallback(agentClient.Signers)
}