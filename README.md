# HA Command to MQTT

[![CI](https://github.com/yourusername/ha-command-to-mqtt/workflows/CI/badge.svg)](https://github.com/yourusername/ha-command-to-mqtt/actions/workflows/ci.yml)
[![Quality](https://github.com/yourusername/ha-command-to-mqtt/workflows/Quality/badge.svg)](https://github.com/yourusername/ha-command-to-mqtt/actions/workflows/quality.yml)
[![Release](https://github.com/yourusername/ha-command-to-mqtt/workflows/Release%20Management/badge.svg)](https://github.com/yourusername/ha-command-to-mqtt/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/ha-command-to-mqtt)](https://goreportcard.com/report/github.com/yourusername/ha-command-to-mqtt)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A Go application that executes system commands and publishes their results to MQTT topics with Home Assistant auto-discovery support.

## Features

- Execute custom commands at configurable intervals
- Publish command results to MQTT topics
- Home Assistant auto-discovery for automatic sensor creation
- Support for both YAML configuration and environment variables
- Configurable device classes, units, and icons for sensors
- Advanced logging with multiple formats (text, JSON, logfmt) using logrus
- SSH support for remote command execution
- Robust error handling and comprehensive logging

## Configuration

### Option 1: YAML Configuration (Recommended)

Create a `config.yaml` file:

```yaml
mqtt:
  broker: "localhost"
  port: 1883
  username: "your_username"
  password: "your_password"
  client_id: "ha-command-to-mqtt"

commands:
  - name: "CPU Temperature"
    command: "cat /sys/class/thermal/thermal_zone0/temp"
    frequency: "30s"
    device_class: "temperature"
    unit: "°C"
    icon: "mdi:thermometer"

  - name: "System Uptime"
    command: "uptime -p"
    frequency: "5m"
    icon: "mdi:clock-outline"
```

### Option 2: Environment Variables

Set environment variables for configuration:

```bash
# MQTT Configuration
export MQTT_BROKER=localhost
export MQTT_PORT=1883
export MQTT_USERNAME=your_username
export MQTT_PASSWORD=your_password
export MQTT_CLIENT_ID=ha-command-to-mqtt

# Commands
export COMMAND_CPU_TEMP="cat /sys/class/thermal/thermal_zone0/temp"
export COMMAND_CPU_TEMP_FREQUENCY=30s
export COMMAND_CPU_TEMP_DEVICE_CLASS=temperature
export COMMAND_CPU_TEMP_UNIT=°C
export COMMAND_CPU_TEMP_ICON=mdi:thermometer
```

## Command Configuration

Each command supports the following options:

- `name`: Display name for the sensor
- `command`: Shell command to execute
- `frequency`: How often to run the command (e.g., "30s", "5m", "1h")
- `device_class`: Home Assistant device class (optional)
- `unit`: Unit of measurement (optional)
- `icon`: Material Design Icon (optional)
- `target_host`: Target host to execute command on - "local" or SSH host name (optional, defaults to "local")
- `force_update`: Force Home Assistant to update even if value hasn't changed (optional)
- `state_class`: Home Assistant state class - "measurement", "total", "total_increasing" (optional)
- `entity_category`: Home Assistant entity category - "config", "diagnostic" (optional)
- `expire_after`: Seconds after which the sensor becomes unavailable if no update (optional)

## SSH Support (Optional)

Commands can be executed either locally or on remote hosts via SSH. **SSH configuration is completely optional** - the application works perfectly fine with only local commands.

### SSH Configuration

If you want to execute commands on remote hosts, configure SSH hosts in your `config.yaml`:

```yaml
ssh:
  hosts:
    - name: "server1"
      host: "192.168.1.100"
      port: 22
      user: "pi"
      key_path: "/home/user/.ssh/id_rsa"
      timeout: "30s"

    - name: "server2"
      host: "example.com"
      port: 2222
      user: "ubuntu"
      key_path: "/home/user/.ssh/id_ed25519"
      timeout: "10s"

    # SSH Agent authentication (no key_path needed)
    - name: "server3"
      host: "agent-server.local"
      port: 22
      user: "ubuntu"
      timeout: "30s"
```

### SSH Authentication Methods

The application supports multiple SSH authentication methods in the following priority order:

1. **Configured SSH Keys**: If `key_path` is specified in the host configuration, it will use that specific key
2. **SSH Agent**: If SSH agent is running (`ssh-agent`) and has keys loaded, it will use those keys
3. **Default SSH Keys**: Falls back to searching for keys in default locations (`~/.ssh/id_rsa`, `~/.ssh/id_ed25519`, `~/.ssh/id_ecdsa`)
4. **Password Authentication**: If configured in the host settings

#### Using SSH Agent

SSH agent provides the most convenient authentication method as it doesn't require specifying key paths in your configuration:

```bash
# Start SSH agent and add your keys
eval $(ssh-agent -s)
ssh-add ~/.ssh/id_rsa
ssh-add ~/.ssh/id_ed25519

# Verify keys are loaded
ssh-add -l

# Run the application - it will automatically use agent keys
./ha-command-to-mqtt
```

With SSH agent, your host configuration becomes simpler:

```yaml
ssh:
  hosts:
    - name: "my-server"
      host: "example.com"
      user: "ubuntu"
      # No key_path needed - will use SSH agent automatically
```

**Alternative SSH configurations:**

```yaml
# Option 1: No SSH section at all (local commands only)
mqtt:
  broker: "localhost"
  # ... other MQTT config

commands:
  - name: "Local Command"
    command: "uptime -p"
    frequency: "5m"

# Option 2: Empty SSH section (same as option 1)
mqtt:
  broker: "localhost"
  # ... other MQTT config

ssh:
  hosts: []

commands:
  - name: "Local Command"
    command: "uptime -p"
    frequency: "5m"
```

### Local and Remote Commands

Commands execute locally by default. To execute on a remote host, specify the `target_host` parameter:

```yaml
commands:
  # Local command (default behavior)
  - name: "Local CPU Usage"
    command: "top -bn1 | grep 'Cpu(s)' | awk '{print $2}'"
    frequency: "30s"

  # Explicit local command (same as above)
  - name: "Local Memory Usage"
    command: "free -m | awk 'NR==2{printf \"%.1f\", $3*100/$2}'"
    frequency: "30s"
    target_host: "local"

  # Remote command via SSH
  - name: "Remote Server Load"
    command: "uptime | awk -F'load average:' '{print $2}' | awk '{print $1}'"
    frequency: "1m"
    target_host: "server1"
    icon: "mdi:server"
```

### SSH Features

- **SSH Agent Support**: Automatically uses keys loaded in ssh-agent for seamless authentication
- **Automatic Key Detection**: If no key path is specified, the application will try SSH agent first, then common SSH key locations (`~/.ssh/id_rsa`, `~/.ssh/id_ed25519`, `~/.ssh/id_ecdsa`)
- **Connection Pooling**: SSH connections are established once and reused for multiple commands
- **Automatic Reconnection**: If an SSH connection drops, the application will automatically reconnect
- **Timeout Support**: Configure connection timeouts per host
- **Multiple Authentication**: Supports SSH agent, SSH key files, and password authentication

### Environment Variables for Remote Commands

```bash
# SSH hosts are configured in YAML, but commands can reference them
COMMAND_REMOTE_CPU=cat /proc/loadavg | awk '{print $1}'
COMMAND_REMOTE_CPU_TARGET_HOST=server1
COMMAND_REMOTE_CPU_FREQUENCY=1m
COMMAND_REMOTE_CPU_STATE_CLASS=measurement

# Local command (optional to specify)
COMMAND_LOCAL_DISK=df -h / | awk 'NR==2{print $5}' | sed 's/%//'
COMMAND_LOCAL_DISK_TARGET_HOST=local
COMMAND_LOCAL_DISK_FREQUENCY=5m
COMMAND_LOCAL_DISK_ENTITY_CATEGORY=diagnostic
COMMAND_LOCAL_DISK_FORCE_UPDATE=true
```

## Home Assistant Integration Attributes

The application supports comprehensive Home Assistant sensor attributes:

### State Classes
- `measurement`: For sensors that report current values (temperature, CPU usage)
- `total`: For sensors that report cumulative values (disk space used)
- `total_increasing`: For monotonically increasing values (network bytes sent)

### Entity Categories
- `config`: Configuration entities
- `diagnostic`: Diagnostic entities (shown in diagnostic section)

### Device Classes
Common device classes include: `battery`, `temperature`, `humidity`, `pressure`, `signal_strength`, `data_size`, `data_rate`, `duration`, `frequency`, `power`, `voltage`, `current`, `energy`, etc.

### Additional Attributes
- `force_update`: Forces Home Assistant to update even if the value hasn't changed
- `expire_after`: Number of seconds after which the sensor becomes unavailable if no new data

Example with all attributes:
```yaml
- name: "Advanced CPU Sensor"
  command: "top -bn1 | grep 'Cpu(s)' | awk '{print $2}' | sed 's/%us,//'"
  frequency: "30s"
  device_class: "battery"
  unit: "%"
  icon: "mdi:cpu-64-bit"
  state_class: "measurement"
  entity_category: "diagnostic"
  force_update: true
  expire_after: 90
```

## Installation

### Home Assistant Add-on (Recommended)

The easiest way to use this application is as a Home Assistant add-on:

1. **Build and install the add-on**:
   ```bash
   # Clone the repository
   git clone https://github.com/yourusername/ha-command-to-mqtt.git
   cd ha-command-to-mqtt

   # Build the add-on
   ./addon/build.sh

   # Copy to Home Assistant addons directory
   scp -r addon/ root@your-ha-ip:/addons/command-to-mqtt/
   ```

2. **Install through Home Assistant UI**:
   - Go to **Supervisor** → **Add-on Store** → **Local add-ons**
   - Find "Command to MQTT" and click **Install**
   - Configure through the add-on configuration UI

📖 **See [addon/INSTALL.md](addon/INSTALL.md) for detailed installation instructions.**

### Manual Installation

1. **Clone or download the project**

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Configure the application** using either `config.yaml` or environment variables

4. **Run the application**:
   ```bash
   go run main.go
   ```

5. **Build for production**:
   ```bash
   go build -o ha-command-to-mqtt
   ./ha-command-to-mqtt
   ```

## Command Line Options

The application supports several command-line flags:

```bash
./ha-command-to-mqtt [OPTIONS]
```

### Available Options

- `-c, --config FILE`: Configuration file path (default: `config.yaml`)
- `-l, --log-level LEVEL`: Log level - `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` (default: `info`)
- `-f, --log-format FORMAT`: Log format - `text`, `json`, `logfmt` (default: `text`)
- `-v, --version`: Show version information and exit
- `-h, --help`: Show help message and exit

### Examples

```bash
# Use default config.yaml and info log level
./ha-command-to-mqtt

# Specify custom config file
./ha-command-to-mqtt --config /path/to/myconfig.yaml

# Set debug log level with JSON format
./ha-command-to-mqtt --log-level debug --log-format json

# Use logfmt format for structured logging
./ha-command-to-mqtt --log-format logfmt

# Combine options using short flags
./ha-command-to-mqtt -c /etc/ha-mqtt/config.yaml -l warn -f json

# Using long form with equals
./ha-command-to-mqtt --config=/path/to/config.yaml --log-level=debug --log-format=logfmt
```

### Environment Variables for CLI Options

You can also set these via environment variables:

```bash
export HA_MQTT_CONFIG="/path/to/config.yaml"
export HA_MQTT_LOG_LEVEL="debug"
export HA_MQTT_LOG_FORMAT="json"
./ha-command-to-mqtt
```

## Logging

The application uses [logrus](https://github.com/sirupsen/logrus) for advanced logging capabilities.

### Log Levels

- `panic`: Highest level, logs and then calls `panic()`
- `fatal`: Logs and then calls `os.Exit(1)`
- `error`: Error messages only
- `warn`: Warning and error messages
- `info`: Standard informational logging (default)
- `debug`: Verbose logging with file names and line numbers
- `trace`: Most verbose level for detailed debugging

### Log Formats

- `text`: Human-readable text format with colors and timestamps (default)
- `json`: Structured JSON format, ideal for log aggregation systems
- `logfmt`: Key-value pairs format (key=value), uses logrus TextFormatter without colors

### Examples

```bash
# Debug level with JSON format (great for log analysis tools)
./ha-command-to-mqtt --log-level debug --log-format json

# Production setup with structured logging
./ha-command-to-mqtt --log-level info --log-format logfmt

# Development with detailed text logging
./ha-command-to-mqtt --log-level trace --log-format text
```

## Home Assistant Integration

The application automatically creates sensors in Home Assistant through MQTT discovery. Sensors will appear in Home Assistant with the device name "Command Sensors".

### MQTT Topics

- Discovery: `homeassistant/sensor/{client_id}_{sensor_name}/config`
- State: `homeassistant/sensor/{client_id}_{sensor_name}/state`

## Example Commands

### System Monitoring Commands

```yaml
commands:
  # CPU Temperature (Linux) - with measurement state class and expiry
  - name: "CPU Temperature"
    command: "cat /sys/class/thermal/thermal_zone0/temp"
    frequency: "30s"
    device_class: "temperature"
    unit: "°C"
    state_class: "measurement"
    expire_after: 120
    icon: "mdi:thermometer"

  # Memory Usage - with force update and measurement class
  - name: "Memory Usage"
    command: "free | grep Mem | awk '{printf \"%.1f\", $3/$2 * 100.0}'"
    frequency: "1m"
    unit: "%"
    device_class: "battery"
    state_class: "measurement"
    force_update: true
    icon: "mdi:memory"

  # Disk Usage - with diagnostic category
  - name: "Disk Usage Root"
    command: "df -h / | awk 'NR==2{print $5}' | sed 's/%//'"
    frequency: "10m"
    unit: "%"
    state_class: "measurement"
    entity_category: "diagnostic"
    icon: "mdi:harddisk"

  # Network Interface Status - simple status sensor
  - name: "Network Interface"
    command: "cat /sys/class/net/eth0/operstate"
    frequency: "1m"
    entity_category: "diagnostic"
    icon: "mdi:ethernet"

  # Docker Container Count - total increasing counter
  - name: "Docker Containers"
    command: "docker ps -q | wc -l"
    frequency: "2m"
    state_class: "total"
    icon: "mdi:docker"
```

### macOS Commands

```yaml
commands:
  # CPU Temperature (requires additional tools)
  - name: "CPU Temperature"
    command: "sudo powermetrics -n 1 -s smc | grep 'CPU die temperature' | awk '{print $4}'"
    frequency: "30s"
    unit: "°C"

  # Memory Pressure
  - name: "Memory Pressure"
    command: "memory_pressure | grep 'System-wide memory free percentage' | awk '{print $5}' | sed 's/%//'"
    frequency: "1m"
    unit: "%"

  # Battery Level (MacBook)
  - name: "Battery Level"
    command: "pmset -g batt | grep -Eo '[0-9]+%' | sed 's/%//'"
    frequency: "5m"
    unit: "%"
    device_class: "battery"
```

### Remote SSH Commands

```yaml
commands:
  # Remote Raspberry Pi monitoring
  - name: "Pi CPU Temperature"
    command: "vcgencmd measure_temp | sed 's/temp=//' | sed 's/°C//'"
    frequency: "1m"
    target_host: "raspberry-pi"
    device_class: "temperature"
    unit: "°C"
    icon: "mdi:thermometer"

  # Remote server disk usage
  - name: "Server Root Disk"
    command: "df -h / | tail -1 | awk '{print $5}' | sed 's/%//'"
    frequency: "5m"
    target_host: "web-server"
    unit: "%"
    icon: "mdi:harddisk"

  # Remote Docker container status
  - name: "Remote Docker Containers"
    command: "docker ps -q | wc -l"
    frequency: "2m"
    target_host: "docker-host"
    icon: "mdi:docker"

  # Remote service status
  - name: "Nginx Status"
    command: "systemctl is-active nginx"
    frequency: "1m"
    target_host: "web-server"
    icon: "mdi:web"

  # Remote network interface
  - name: "Server Network TX"
    command: "cat /sys/class/net/eth0/statistics/tx_bytes"
    frequency: "30s"
    target_host: "gateway"
    icon: "mdi:ethernet"
```

## Docker Support

You can also run this in a Docker container. Here's a basic Dockerfile:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o ha-command-to-mqtt

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ha-command-to-mqtt .
CMD ["./ha-command-to-mqtt"]
```

## Systemd Service

Create a systemd service file at `/etc/systemd/system/ha-command-to-mqtt.service`:

```ini
[Unit]
Description=HA Command to MQTT
After=network.target

[Service]
Type=simple
User=your_user
WorkingDirectory=/path/to/ha-command-to-mqtt
ExecStart=/path/to/ha-command-to-mqtt/ha-command-to-mqtt
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl enable ha-command-to-mqtt
sudo systemctl start ha-command-to-mqtt
```

## Development and CI/CD

This project includes a comprehensive CI/CD pipeline system using GitHub Actions.

### Automated Workflows

- **CI Pipeline** (`ci.yml`): Automated testing, building, and deployment
  - Go testing with multiple platforms (Linux, macOS, Windows)
  - Multi-architecture Docker image builds
  - Home Assistant add-on building and testing
  - Automated releases with GitHub Releases

- **Quality Assurance** (`quality.yml`): Code quality and security scanning
  - `golangci-lint` for Go code quality
  - `gosec` for security vulnerability scanning
  - Trivy for Docker image vulnerability scanning
  - CodeQL for advanced code analysis
  - Dependency vulnerability monitoring

- **Release Management** (`release.yml`): Automated release creation
  - Multi-platform binary builds (Linux, macOS, Windows)
  - Home Assistant add-on package creation
  - Automated changelog generation
  - Release asset creation and publishing

- **Dependency Management** (`dependencies.yml`): Automated dependency updates
  - Go module dependency updates
  - GitHub Actions version updates
  - Security vulnerability monitoring
  - Base Docker image update notifications

### Development Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/ha-command-to-mqtt.git
   cd ha-command-to-mqtt
   ```

2. **Install development dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run quality checks**:
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Run linting
   golangci-lint run
   
   # Run tests
   go test ./...
   ```

4. **Build locally**:
   ```bash
   go build -o ha-command-to-mqtt
   ```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run quality checks locally
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

All pull requests are automatically tested through the CI pipeline.

## Troubleshooting

- Check that MQTT broker is accessible
- Verify command syntax in your shell before adding to config
- Monitor logs for command execution errors
- Ensure Home Assistant has MQTT discovery enabled
- Check MQTT topic structure matches Home Assistant expectations
- Review GitHub Actions logs for build and deployment issues

## License

This project is open source and available under the MIT License.