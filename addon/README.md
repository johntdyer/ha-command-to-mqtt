# Command to MQTT - Home Assistant Add-on

This directory contains the Home Assistant add-on configuration for the Command to MQTT application.

## Files Overview

- **`config.yaml`**: Home Assistant add-on configuration and schema
- **`DOCS.md`**: Add-on documentation shown in Home Assistant UI
- **`Dockerfile`**: Container build instructions for the add-on
- **`run.sh`**: Entry point script that converts HA options to application config
- **`build.sh`**: Build script for creating multi-architecture binaries
- **`README.md`**: This file

## Installation

### Method 1: Local Add-on Repository

1. **Build the add-on**:
   ```bash
   cd /path/to/ha-command-to-mqtt
   ./addon/build.sh
   ```

2. **Copy to Home Assistant**:
   ```bash
   # Copy the entire addon directory to your Home Assistant addons folder
   cp -r addon /path/to/homeassistant/addons/command-to-mqtt
   ```

3. **Install through Home Assistant**:
   - Go to **Supervisor** → **Add-on Store** → **Local add-ons**
   - Find "Command to MQTT" and click **Install**

### Method 2: Custom Repository

1. **Create a repository** (e.g., on GitHub) with this structure:
   ```
   my-ha-addons/
   ├── repository.yaml
   └── command-to-mqtt/
       ├── config.yaml
       ├── DOCS.md
       ├── Dockerfile
       └── run.sh
   ```

2. **Add repository to Home Assistant**:
   - Go to **Supervisor** → **Add-on Store** → **⋮** → **Repositories**
   - Add your repository URL
   - Install the "Command to MQTT" add-on

## Configuration

The add-on configuration is done through the Home Assistant UI. Here's a typical configuration:

```yaml
mqtt:
  broker: "core-mosquitto"
  port: 1883
  username: "your_mqtt_user"
  password: "your_mqtt_password"
  client_id: "ha-command-to-mqtt"
  discovery_prefix: "homeassistant"

commands:
  - name: "system_uptime"
    command: "uptime"
    topic: "system/uptime"
    interval: "300s"
    unit_of_measurement: ""
    device_class: "timestamp"
    
  - name: "cpu_usage"
    command: "top -bn1 | grep 'Cpu(s)' | awk '{print $2}' | sed 's/%us,//'"
    topic: "system/cpu_usage"
    interval: "60s"
    unit_of_measurement: "%"
    state_class: "measurement"

ssh:
  hosts:
    - name: "raspberry-pi"
      host: "192.168.1.100"
      port: 22
      user: "pi"
      key_path: "/share/ssh_keys/id_rsa"
      timeout: "30s"

log_level: "info"
log_format: "text"
```

## SSH Key Management

### Using Shared Directory

1. **Place SSH keys in the shared directory**:
   ```bash
   # On Home Assistant host
   cp ~/.ssh/id_rsa /usr/share/hassio/share/ssh_keys/
   chmod 600 /usr/share/hassio/share/ssh_keys/id_rsa
   ```

2. **Reference in configuration**:
   ```yaml
   ssh:
     hosts:
       - name: "my-server"
         host: "example.com"
         user: "ubuntu"
         key_path: "/share/ssh_keys/id_rsa"
   ```

### SSH Agent Support

The add-on automatically detects SSH agent if available. This is useful for advanced setups where SSH agent is running on the Home Assistant host.

## Building Multi-Architecture Images

For production deployment, build for all supported architectures:

```bash
# Build for all Home Assistant supported architectures
./addon/build.sh --all-arch

# This creates binaries for:
# - linux/amd64 (Intel/AMD 64-bit)
# - linux/arm64 (ARM 64-bit, like Raspberry Pi 4)
# - linux/arm/v7 (ARM 32-bit v7, like Raspberry Pi 3)
# - linux/arm/v6 (ARM 32-bit v6, like Raspberry Pi Zero)
# - linux/386 (Intel/AMD 32-bit)
```

## Development

### Testing Locally

1. **Build the add-on**:
   ```bash
   ./addon/build.sh
   ```

2. **Test the generated config conversion**:
   ```bash
   # The run.sh script converts Home Assistant options to YAML config
   # Test this locally by setting bashio environment
   ```

### Debugging

- Check add-on logs in Home Assistant: **Supervisor** → **Command to MQTT** → **Log**
- Enable debug logging by setting `log_level: "debug"` in the add-on configuration
- Use `log_format: "json"` for structured logging

## Supported Home Assistant Features

- ✅ **MQTT Discovery**: Automatic sensor creation
- ✅ **Add-on Configuration**: UI-based configuration
- ✅ **Multi-Architecture**: Supports all HA architectures
- ✅ **Shared Storage**: Access to `/share` directory
- ✅ **SSL Certificates**: Access to `/ssl` directory (read-only)
- ✅ **Configuration Access**: Access to `/config` directory (read-only)
- ✅ **Supervisor API**: Integration with Home Assistant supervisor
- ✅ **Auto-start**: Automatic startup with Home Assistant

## Troubleshooting

### Add-on Won't Start
- Check the add-on logs for error messages
- Verify MQTT broker configuration
- Ensure required fields are filled in configuration

### SSH Connection Issues
- Verify SSH key permissions (should be 600)
- Check if the target host is reachable
- Test SSH connection manually from Home Assistant host

### MQTT Issues
- Verify MQTT broker is running (`core-mosquitto` add-on)
- Check MQTT credentials
- Ensure MQTT discovery is enabled in Home Assistant

### Configuration Errors
- Validate YAML syntax in add-on configuration
- Check that all required fields are provided
- Review the add-on logs for configuration parsing errors

## Support

For issues specific to the Home Assistant add-on, please check:
1. Add-on logs in Home Assistant UI
2. Home Assistant supervisor logs
3. The main project repository for application-specific issues