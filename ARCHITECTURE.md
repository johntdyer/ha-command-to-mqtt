# Project Architecture and Structure

## Overview

HA Command to MQTT is a comprehensive Go application designed for executing system commands and publishing results to MQTT topics with Home Assistant auto-discovery. The project has evolved from a simple command executor to a full enterprise-grade solution with modular architecture, SSH support, Home Assistant integration, and complete CI/CD automation.

## Project Structure

```
ha-command-to-mqtt/
├── .github/                     # GitHub Actions workflows and templates
│   ├── workflows/
│   │   ├── ci.yml              # Main CI/CD pipeline
│   │   ├── addon-release.yml   # Home Assistant add-on releases
│   │   ├── quality.yml         # Code quality and security scanning
│   │   ├── dependencies.yml    # Dependency management automation
│   │   └── release.yml         # Release management
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.yml      # Bug report template
│   │   ├── feature_request.yml # Feature request template
│   │   └── addon_issue.yml     # Add-on specific issues
│   └── pull_request_template.md # PR template
├── addon/                       # Home Assistant add-on files
│   ├── config.yaml             # Add-on configuration
│   ├── Dockerfile              # Add-on container build
│   ├── run.sh                  # Add-on startup script
│   ├── build.sh               # Add-on build script
│   ├── DOCS.md                # Add-on documentation
│   └── icon.svg               # Add-on icon
├── docker/                      # Docker deployment files
│   ├── Dockerfile              # Application container
│   ├── docker-compose.yml      # Docker Compose setup
│   └── README.md               # Docker documentation
├── config/                      # Configuration management
│   └── config.go              # Configuration loading and validation
├── internal/                    # Internal application modules
│   ├── cli/                   # Command line interface
│   │   └── cli.go
│   ├── executor/              # Command execution engine
│   │   └── executor.go
│   ├── logging/               # Advanced logging system
│   │   └── logging.go
│   ├── mqtt/                  # MQTT client and Home Assistant integration
│   │   └── mqtt.go
│   └── ssh/                   # SSH client and authentication
│       └── ssh.go
├── .golangci.yml               # Go linting configuration
├── .yamllint.yml               # YAML linting configuration
├── config.yaml                 # Application configuration example
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── main.go                     # Application entry point
└── README.md                   # Project documentation
```

## Architecture Components

### Core Application (`main.go`)

The main entry point orchestrates all components:
- CLI argument parsing and configuration loading
- Logger initialization with multiple formats
- SSH client pool management
- MQTT client setup and connection management
- Command execution scheduling and coordination
- Graceful shutdown handling

### Configuration Management (`config/config.go`)

Handles application configuration through multiple sources:
- YAML configuration files with full validation
- Environment variable override support
- Command-line flag integration via Viper
- Configuration validation and defaults
- SSH host configuration management

### Command Line Interface (`internal/cli/cli.go`)

Professional CLI using Cobra and Viper:
- Command-line flag definitions and parsing
- Configuration file path specification
- Log level and format control
- Version information display
- Help system integration

### Logging System (`internal/logging/logging.go`)

Advanced logging using Logrus:
- Multiple output formats: text, JSON, logfmt
- Configurable log levels (panic to trace)
- Structured logging with context fields
- File and console output support
- Production-ready log formatting

### SSH Client (`internal/ssh/ssh.go`)

Comprehensive SSH support with multiple authentication methods:
- SSH agent integration for seamless authentication
- Private key file support (RSA, Ed25519, ECDSA)
- Connection pooling and automatic reconnection
- Timeout configuration per host
- Command execution over SSH with proper output handling
- Error handling and connection state management

### MQTT Client (`internal/mqtt/mqtt.go`)

Home Assistant integration via MQTT:
- Eclipse Paho MQTT client integration
- Home Assistant auto-discovery protocol implementation
- Sensor configuration with full attribute support
- Topic management and message publishing
- Connection state monitoring and automatic reconnection
- Device class, state class, and entity category support

### Command Executor (`internal/executor/executor.go`)

Command execution engine with scheduling:
- Local command execution with proper output handling
- Remote command execution via SSH
- Configurable execution frequency per command
- Multi-line output support with proper parsing
- Error handling and logging integration
- Command result processing and MQTT publishing

## Home Assistant Add-on (`addon/`)

Complete Home Assistant Supervisor add-on:
- Native Home Assistant integration
- UI-based configuration management
- Multi-architecture Docker support (amd64, arm64, armhf, armv7, i386)
- Automatic service discovery and setup
- Add-on store compatibility
- Comprehensive documentation and installation guide

## CI/CD Pipeline System (`.github/workflows/`)

Comprehensive automation covering the entire development lifecycle:

### Main CI Pipeline (`ci.yml`)
- **Testing**: Go vet, staticcheck, unit tests across multiple platforms
- **Building**: Multi-platform binary builds (Linux, macOS, Windows)
- **Docker**: Multi-architecture container images with proper tagging
- **Add-on**: Home Assistant add-on building and testing
- **Security**: Vulnerability scanning with Trivy
- **Releases**: Automated GitHub releases with proper versioning

### Quality Assurance (`quality.yml`)
- **Code Quality**: golangci-lint with comprehensive rule set
- **Security Scanning**: gosec for Go security vulnerabilities
- **Container Security**: Trivy scanning for Docker images
- **Advanced Analysis**: CodeQL for deep code analysis
- **Dependency Security**: Vulnerability monitoring for all dependencies
- **License Compliance**: License compatibility checking

### Release Management (`release.yml`)
- **Multi-platform Releases**: Automated binary builds for all platforms
- **Add-on Packaging**: Complete Home Assistant add-on packages
- **Changelog Generation**: Automated changelog from git history
- **Asset Management**: Checksums and release asset organization
- **Version Management**: Semantic versioning and tag management

### Dependency Management (`dependencies.yml`)
- **Automated Updates**: Go module and GitHub Actions updates
- **Security Monitoring**: Continuous vulnerability scanning
- **Update Notifications**: Automated dependency update PRs
- **Base Image Updates**: Docker base image update reminders

## Key Features

### Modularity
- Clean separation of concerns across focused modules
- Interface-based design for easy testing and extension
- Dependency injection pattern for component integration
- Configuration-driven behavior modification

### Scalability
- Connection pooling for SSH clients
- Concurrent command execution with proper synchronization
- Configurable execution frequencies per command
- Resource-efficient MQTT client management

### Reliability
- Comprehensive error handling and recovery
- Automatic reconnection for SSH and MQTT
- Graceful shutdown with proper cleanup
- Extensive logging for troubleshooting

### Security
- Multiple SSH authentication methods
- Secure credential handling
- No hardcoded secrets or credentials
- Security scanning integration in CI/CD

### Observability
- Structured logging with multiple formats
- Comprehensive error reporting
- Performance monitoring capabilities
- Debug-level tracing support

## Development Workflow

### Local Development
1. **Setup**: Clone repository and install dependencies
2. **Configuration**: Create local config.yaml for testing
3. **Development**: Make changes with proper module organization
4. **Testing**: Run local tests and quality checks
5. **Validation**: Test with real MQTT broker and Home Assistant

### CI/CD Integration
1. **Commit**: Push changes trigger automated testing
2. **Quality**: Automated code quality and security scanning
3. **Building**: Multi-platform builds and Docker images
4. **Testing**: Integration testing with Home Assistant add-on
5. **Release**: Automated releases with proper versioning

### Production Deployment
1. **Home Assistant Add-on**: Native HA integration with UI configuration
2. **Docker Deployment**: Container-based deployment with compose
3. **Binary Installation**: Direct binary installation with systemd
4. **SSH Configuration**: Flexible SSH authentication setup

## Future Extensibility

The modular architecture supports easy extension:
- **New Command Types**: Easy addition of specialized command executors
- **Additional Protocols**: Support for other messaging protocols
- **Enhanced Authentication**: Additional SSH authentication methods
- **Monitoring Integration**: Prometheus metrics or other monitoring systems
- **Configuration Sources**: Additional configuration backends (Consul, etcd)
- **Output Formats**: Additional Home Assistant sensor types and formats

## Technology Stack

- **Core Language**: Go 1.21+ with modern language features
- **MQTT Client**: Eclipse Paho MQTT for reliable messaging
- **SSH Support**: golang.org/x/crypto/ssh with agent integration
- **Configuration**: Viper for flexible configuration management
- **CLI Framework**: Cobra for professional command-line interface
- **Logging**: Logrus for production-grade logging
- **Containerization**: Docker with multi-architecture support
- **CI/CD**: GitHub Actions with comprehensive automation
- **Quality Tools**: golangci-lint, gosec, Trivy, CodeQL
- **Home Assistant**: Native add-on integration with Supervisor

This architecture provides a solid foundation for a production-ready application with enterprise-grade features, comprehensive automation, and excellent maintainability.