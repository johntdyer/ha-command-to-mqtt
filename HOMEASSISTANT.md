# Home Assistant MQTT Discovery Message Example

This shows what the application generates for Home Assistant auto-discovery.

## Example Discovery Message

For a command configured as:
```yaml
- name: "CPU Temperature"
  command: "cat /sys/class/thermal/thermal_zone0/temp"
  frequency: "30s"
  device_class: "temperature"
  unit: "°C"
  icon: "mdi:thermometer"
  state_class: "measurement"
  force_update: true
  expire_after: 120
```

The application publishes this discovery message to:
`homeassistant/sensor/ha-command-to-mqtt_cpu_temperature/config`

```json
{
  "name": "CPU Temperature",
  "state_topic": "homeassistant/sensor/ha-command-to-mqtt_cpu_temperature/state",
  "unique_id": "ha-command-to-mqtt_cpu_temperature",
  "device_class": "temperature",
  "unit_of_measurement": "°C",
  "icon": "mdi:thermometer",
  "state_class": "measurement",
  "force_update": true,
  "expire_after": 120,
  "device": {
    "identifiers": ["ha-command-to-mqtt"],
    "name": "Command Sensors",
    "model": "HA Command to MQTT",
    "manufacturer": "Custom"
  }
}
```

And publishes the actual temperature value to:
`homeassistant/sensor/ha-command-to-mqtt_cpu_temperature/state`

## Supported Home Assistant Attributes

### State Classes
- `measurement`: Current readings (temperature, CPU %, etc.)
- `total`: Cumulative values that may decrease (disk usage)
- `total_increasing`: Monotonically increasing values (network counters)

### Device Classes
- `battery`: Percentage values (0-100%)
- `temperature`: Temperature readings
- `humidity`: Humidity percentage
- `pressure`: Pressure readings
- `data_size`: File/memory sizes
- `data_rate`: Transfer rates
- `power`: Power consumption
- `energy`: Energy consumption
- `voltage`: Electrical voltage
- `current`: Electrical current
- And many more...

### Entity Categories
- `config`: Configuration entities
- `diagnostic`: Diagnostic information (shown in separate section)

### Other Attributes
- `force_update`: Always trigger automation even if value unchanged
- `expire_after`: Seconds until sensor becomes unavailable without updates
- `icon`: Material Design Icon (mdi:icon-name)
- `unit_of_measurement`: Units for the sensor value