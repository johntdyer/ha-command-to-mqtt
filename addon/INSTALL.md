# Home Assistant Add-on Installation Guide

This guide will help you install and configure the Command to MQTT add-on in Home Assistant.

## Prerequisites

- Home Assistant Supervisor (formerly Hass.io)
- MQTT integration configured in Home Assistant
- Basic understanding of YAML configuration

## Installation Methods

### Method 1: Local Add-on (Recommended for Testing)

1. **Prepare your Home Assistant system**:
   - Ensure you have SSH access to your Home Assistant host
   - Make sure you have the MQTT broker add-on installed and configured

2. **Build and copy the add-on**:
   ```bash
   # On your development machine
   git clone <your-repository>
   cd ha-command-to-mqtt
   
   # Build the add-on
   ./addon/build.sh
   
   # Copy to Home Assistant (adjust paths as needed)
   scp -r addon/ root@your-ha-ip:/addons/command-to-mqtt/
   ```

3. **Install through Home Assistant**:
   - Go to **Supervisor** → **Add-on Store**
   - Refresh the page (Ctrl+F5)
   - Look for "Command to MQTT" in **Local add-ons**
   - Click **Install**

### Method 2: Custom Repository

1. **Create a GitHub repository** with this structure:
   ```
   your-ha-addons/
   ├── repository.yaml
   └── command-to-mqtt/
       ├── config.yaml
       ├── DOCS.md
       ├── Dockerfile
       ├── run.sh
       └── icon.svg
   ```

2. **Configure repository.yaml**:
   ```yaml
   name: "Your Custom Add-ons"
   url: "https://github.com/yourusername/your-ha-addons"
   maintainer: "Your Name <your.email@example.com>"
   ```

3. **Add to Home Assistant**:
   - Go to **Supervisor** → **Add-on Store** → **⋮** (three dots) → **Repositories**
   - Add: `https://github.com/yourusername/your-ha-addons`
   - Install "Command to MQTT" from your custom repository

## Configuration

### Basic Setup

1. **Install and configure MQTT broker** (if not already done):
   - Install the "Mosquitto broker" add-on from the official add-on store
   - Configure with username/password
   - Enable MQTT integration in Home Assistant

2. **Configure the Command to MQTT add-on**:

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
       
     - name: "cpu_temperature"
       command: "cat /sys/class/thermal/thermal_zone0/temp | awk '{print $1/1000}'"
       topic: "system/cpu_temp"
       interval: "60s"
       unit_of_measurement: "°C"
       device_class: "temperature"
       state_class: "measurement"
       
     - name: "memory_usage"
       command: "free | grep Mem | awk '{printf \"%.1f\", $3/$2 * 100.0}'"
       topic: "system/memory_usage"
       interval: "120s"
       unit_of_measurement: "%"
       state_class: "measurement"

   log_level: "info"
   log_format: "text"
   ```

### SSH Configuration

For monitoring remote systems:

1. **Prepare SSH keys**:
   ```bash
   # Generate SSH key pair if you don't have one
   ssh-keygen -t ed25519 -f ~/.ssh/ha_command_key
   
   # Copy public key to target systems
   ssh-copy-id -i ~/.ssh/ha_command_key.pub user@target-host
   
   # Copy private key to Home Assistant shared directory
   scp ~/.ssh/ha_command_key root@ha-ip:/usr/share/hassio/share/ssh_keys/
   ```

2. **Configure SSH hosts**:
   ```yaml
   ssh:
     hosts:
       - name: "raspberry-pi"
         host: "192.168.1.100"
         port: 22
         user: "pi"
         key_path: "/share/ssh_keys/ha_command_key"
         timeout: "30s"
         
       - name: "linux-server"
         host: "192.168.1.101"
         port: 22
         user: "ubuntu"
         key_path: "/share/ssh_keys/ha_command_key"
         timeout: "30s"

   commands:
     - name: "pi_temperature"
       command: "vcgencmd measure_temp | cut -d= -f2 | cut -d\\' -f1"
       topic: "pi/temperature"
       interval: "60s"
       target_host: "raspberry-pi"
       unit_of_measurement: "°C"
       device_class: "temperature"
       state_class: "measurement"
       
     - name: "server_load"
       command: "uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | sed 's/,//'"
       topic: "server/load_average"
       interval: "120s"
       target_host: "linux-server"
       state_class: "measurement"
   ```

## Advanced Configuration Examples

### Network Monitoring

```yaml
commands:
  - name: "internet_speed_down"
    command: "speedtest --simple | grep Download | awk '{print $2}'"
    topic: "network/download_speed"
    interval: "1800s"  # Every 30 minutes
    unit_of_measurement: "Mbps"
    state_class: "measurement"
    
  - name: "ping_response"
    command: "ping -c 1 8.8.8.8 | tail -1| awk '{print $4}' | cut -d '/' -f 2"
    topic: "network/ping_time"
    interval: "60s"
    unit_of_measurement: "ms"
    state_class: "measurement"
```

### Disk Monitoring

```yaml
commands:
  - name: "disk_usage_root"
    command: "df -h / | tail -1 | awk '{print $5}' | sed 's/%//'"
    topic: "system/disk_usage_root"
    interval: "600s"
    unit_of_measurement: "%"
    state_class: "measurement"
    
  - name: "disk_free_gb"
    command: "df -BG / | tail -1 | awk '{print $4}' | sed 's/G//'"
    topic: "system/disk_free"
    interval: "600s"
    unit_of_measurement: "GB"
    state_class: "measurement"
```

### Service Monitoring

```yaml
commands:
  - name: "nginx_status"
    command: "systemctl is-active nginx"
    topic: "services/nginx"
    interval: "120s"
    
  - name: "docker_containers"
    command: "docker ps -q | wc -l"
    topic: "services/docker_containers"
    interval: "300s"
    state_class: "measurement"
```

## Troubleshooting

### Common Issues

1. **Add-on won't start**:
   - Check the add-on logs for specific error messages
   - Verify MQTT broker configuration and connectivity
   - Ensure all required configuration fields are filled

2. **SSH connections failing**:
   - Verify SSH key permissions: `chmod 600 /share/ssh_keys/your_key`
   - Test SSH connection manually: `ssh -i /share/ssh_keys/your_key user@host`
   - Check network connectivity to target hosts

3. **Commands not executing**:
   - Test commands manually on the target system
   - Check command syntax and escaping
   - Verify the command exists on the target system

4. **Sensors not appearing in Home Assistant**:
   - Check MQTT integration is configured and working
   - Verify MQTT discovery is enabled
   - Check MQTT topics are being published (use MQTT explorer)

### Debug Mode

Enable detailed logging:

```yaml
log_level: "debug"
log_format: "json"
```

Then check the add-on logs for detailed execution information.

### Manual Testing

You can test MQTT connectivity and sensor creation manually:

```bash
# Test MQTT publish
mosquitto_pub -h localhost -u your_user -P your_password \
  -t "homeassistant/sensor/test/config" \
  -m '{"name":"Test Sensor","state_topic":"test/sensor","unit_of_measurement":"test"}'

# Test command execution
ssh -i /share/ssh_keys/your_key user@host 'your_command'
```

## Getting Help

1. **Check add-on logs**: Supervisor → Command to MQTT → Log
2. **Enable debug logging**: Set `log_level: "debug"` in configuration
3. **Test components individually**: MQTT, SSH, commands separately
4. **Review documentation**: Check DOCS.md in the add-on for detailed information

## Next Steps

After successful installation:

1. **Start with simple commands** to verify basic functionality
2. **Add SSH hosts gradually** and test each one
3. **Monitor Home Assistant logs** for any integration issues
4. **Create dashboards** using your new sensors
5. **Set up automations** based on command results

The add-on will create sensors automatically in Home Assistant that you can use in dashboards, automations, and scripts.