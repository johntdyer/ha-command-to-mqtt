# Command to MQTT Home Assistant Add-on

Execute local and remote commands and publish their results to MQTT topics that automatically create Home Assistant sensors through MQTT discovery.

## Features

- **Local Command Execution**: Run commands directly on the Home Assistant host
- **SSH Support**: Execute commands on remote servers with SSH key or agent authentication
- **MQTT Discovery**: Automatically creates Home Assistant sensors
- **Flexible Configuration**: YAML-based configuration with environment variable support
- **Multiple Authentication**: SSH keys, SSH agent, and password authentication
- **Advanced Logging**: Configurable log levels and formats (text, JSON, logfmt)
- **Home Assistant Integration**: Built-in MQTT broker support with sensor attributes

## Configuration

### Basic Configuration

```yaml
mqtt:
  broker: "core-mosquitto"  # Use Home Assistant's built-in MQTT broker
  port: 1883
  username: "your_mqtt_user"
  password: "your_mqtt_password"
  client_id: "ha-command-to-mqtt"

commands:
  - name: "system_uptime"
    command: "uptime"
    topic: "system/uptime"
    interval: "300s"
    unit_of_measurement: ""
    device_class: "timestamp"

  - name: "disk_usage"
    command: "df -h / | tail -1 | awk '{print $5}' | sed 's/%//'"
    topic: "system/disk_usage"
    interval: "600s"
    unit_of_measurement: "%"
    device_class: ""
    state_class: "measurement"
```

### SSH Configuration

For remote command execution:

```yaml
ssh:
  hosts:
    - name: "pi-server"
      host: "192.168.1.100"
      port: 22
      user: "pi"
      key_path: "/share/ssh_keys/id_rsa"
      timeout: "30s"

commands:
  - name: "pi_temperature"
    command: "vcgencmd measure_temp | cut -d= -f2 | cut -d\\' -f1"
    topic: "pi/temperature"
    interval: "60s"
    target_host: "pi-server"
    unit_of_measurement: "°C"
    device_class: "temperature"
    state_class: "measurement"
```

### Advanced Options

```yaml
log_level: "info"        # debug, info, warn, error
log_format: "text"       # text, json, logfmt
```

## SSH Key Management

### Option 1: Use Shared Directory
Place your SSH keys in `/share/ssh_keys/` and reference them in the configuration:

```yaml
ssh:
  hosts:
    - name: "server1"
      host: "example.com"
      user: "ubuntu"
      key_path: "/share/ssh_keys/id_rsa"
```

### Option 2: SSH Agent (Advanced)
If you have SSH agent running on the host, the add-on will automatically detect and use loaded keys.

## Home Assistant Sensor Attributes

The add-on supports all Home Assistant sensor attributes:

- **unit_of_measurement**: Unit for the sensor value (°C, %, MB, etc.)
- **device_class**: Sensor type (temperature, humidity, timestamp, etc.)
- **state_class**: How HA should treat the data (measurement, total, total_increasing)
- **entity_category**: Sensor category (config, diagnostic)
- **force_update**: Force updates even if value hasn't changed

## Example Commands

### System Monitoring
```yaml
commands:
  - name: "cpu_usage"
    command: "top -bn1 | grep 'Cpu(s)' | awk '{print $2}' | sed 's/%us,//'"
    topic: "system/cpu_usage"
    interval: "30s"
    unit_of_measurement: "%"
    device_class: ""
    state_class: "measurement"

  - name: "memory_usage"
    command: "free | grep Mem | awk '{printf \"%.1f\", $3/$2 * 100.0}'"
    topic: "system/memory_usage"
    interval: "60s"
    unit_of_measurement: "%"
    state_class: "measurement"

  - name: "load_average"
    command: "uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | sed 's/,//'"
    topic: "system/load_average"
    interval: "60s"
    state_class: "measurement"
```

### Network Monitoring
```yaml
commands:
  - name: "internet_speed"
    command: "speedtest --simple | grep Download | awk '{print $2}'"
    topic: "network/download_speed"
    interval: "1800s"  # Every 30 minutes
    unit_of_measurement: "Mbps"
    state_class: "measurement"
```

## Troubleshooting

### MQTT Connection Issues
- Ensure the MQTT broker is running and accessible
- Check username/password credentials
- Verify network connectivity

### SSH Connection Issues
- Ensure SSH keys have correct permissions (600)
- Verify host connectivity and SSH service is running
- Check SSH key format and authentication

### Command Execution Issues
- Test commands manually first
- Check command syntax and escaping
- Verify target host has required tools/commands

## Support

For issues and feature requests, please visit the project repository.