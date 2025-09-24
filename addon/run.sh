#!/usr/bin/with-contenv bashio

# Exit on error
set -e

bashio::log.info "Starting Command to MQTT add-on..."

# Read configuration from Home Assistant add-on options
CONFIG_PATH="/tmp/addon-config.yaml"

# Get MQTT configuration
MQTT_BROKER=$(bashio::config 'mqtt.broker')
MQTT_PORT=$(bashio::config 'mqtt.port')
MQTT_USERNAME=$(bashio::config 'mqtt.username')
MQTT_PASSWORD=$(bashio::config 'mqtt.password')
MQTT_CLIENT_ID=$(bashio::config 'mqtt.client_id')
MQTT_DISCOVERY_PREFIX=$(bashio::config 'mqtt.discovery_prefix' 'homeassistant')

# Get logging configuration
LOG_LEVEL=$(bashio::config 'log_level' 'info')
LOG_FORMAT=$(bashio::config 'log_format' 'text')

bashio::log.info "Configuring MQTT connection to ${MQTT_BROKER}:${MQTT_PORT}"

# Create configuration file from add-on options
cat > "${CONFIG_PATH}" << EOF
mqtt:
  broker: "${MQTT_BROKER}"
  port: ${MQTT_PORT}
  username: "${MQTT_USERNAME}"
  password: "${MQTT_PASSWORD}"
  client_id: "${MQTT_CLIENT_ID}"
  discovery_prefix: "${MQTT_DISCOVERY_PREFIX}"

EOF

# Add SSH configuration if present
if bashio::config.exists 'ssh.hosts'; then
    echo "ssh:" >> "${CONFIG_PATH}"
    echo "  hosts:" >> "${CONFIG_PATH}"
    
    # Parse SSH hosts from JSON array
    for host in $(bashio::config 'ssh.hosts | keys[]'); do
        name=$(bashio::config "ssh.hosts[${host}].name")
        hostname=$(bashio::config "ssh.hosts[${host}].host")
        port=$(bashio::config "ssh.hosts[${host}].port" "22")
        user=$(bashio::config "ssh.hosts[${host}].user")
        key_path=$(bashio::config "ssh.hosts[${host}].key_path" "")
        password=$(bashio::config "ssh.hosts[${host}].password" "")
        timeout=$(bashio::config "ssh.hosts[${host}].timeout" "30s")
        
        cat >> "${CONFIG_PATH}" << EOF
    - name: "${name}"
      host: "${hostname}"
      port: ${port}
      user: "${user}"
EOF
        
        if [ -n "${key_path}" ]; then
            echo "      key_path: \"${key_path}\"" >> "${CONFIG_PATH}"
        fi
        
        if [ -n "${password}" ]; then
            echo "      password: \"${password}\"" >> "${CONFIG_PATH}"
        fi
        
        echo "      timeout: \"${timeout}\"" >> "${CONFIG_PATH}"
    done
fi

# Add commands configuration
echo "" >> "${CONFIG_PATH}"
echo "commands:" >> "${CONFIG_PATH}"

for command in $(bashio::config 'commands | keys[]'); do
    name=$(bashio::config "commands[${command}].name")
    cmd=$(bashio::config "commands[${command}].command")
    topic=$(bashio::config "commands[${command}].topic")
    interval=$(bashio::config "commands[${command}].interval")
    target_host=$(bashio::config "commands[${command}].target_host" "local")
    unit_of_measurement=$(bashio::config "commands[${command}].unit_of_measurement" "")
    device_class=$(bashio::config "commands[${command}].device_class" "")
    state_class=$(bashio::config "commands[${command}].state_class" "")
    entity_category=$(bashio::config "commands[${command}].entity_category" "")
    force_update=$(bashio::config "commands[${command}].force_update" "false")
    
    cat >> "${CONFIG_PATH}" << EOF
  - name: "${name}"
    command: "${cmd}"
    topic: "${topic}"
    interval: "${interval}"
    target_host: "${target_host}"
EOF
    
    if [ -n "${unit_of_measurement}" ]; then
        echo "    unit_of_measurement: \"${unit_of_measurement}\"" >> "${CONFIG_PATH}"
    fi
    
    if [ -n "${device_class}" ]; then
        echo "    device_class: \"${device_class}\"" >> "${CONFIG_PATH}"
    fi
    
    if [ -n "${state_class}" ]; then
        echo "    state_class: \"${state_class}\"" >> "${CONFIG_PATH}"
    fi
    
    if [ -n "${entity_category}" ]; then
        echo "    entity_category: \"${entity_category}\"" >> "${CONFIG_PATH}"
    fi
    
    echo "    force_update: ${force_update}" >> "${CONFIG_PATH}"
done

bashio::log.info "Configuration generated successfully"
bashio::log.debug "Generated config:"
if bashio::debug; then
    cat "${CONFIG_PATH}"
fi

# Export environment variables for SSH agent if available
if [ -S "/tmp/ssh-agent/socket" ]; then
    export SSH_AUTH_SOCK="/tmp/ssh-agent/socket"
    bashio::log.info "SSH agent socket detected"
fi

# Start the application with generated configuration
bashio::log.info "Starting ha-command-to-mqtt with log level: ${LOG_LEVEL}, format: ${LOG_FORMAT}"

exec /usr/local/bin/ha-command-to-mqtt \
    --config "${CONFIG_PATH}" \
    --log-level "${LOG_LEVEL}" \
    --log-format "${LOG_FORMAT}"