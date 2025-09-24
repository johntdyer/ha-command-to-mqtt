package main

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

// InitMQTT connects to MQTT broker
func InitMQTT(config *MQTTConfig) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Broker, config.Port))
	opts.SetClientID(config.ClientID)

	if config.Username != "" {
		opts.SetUsername(config.Username)
	}
	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		logger.Debugf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
	})

	mqttClient = mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	logger.Info("Connected to MQTT broker")
	return nil
}

// DisconnectMQTT disconnects from MQTT broker
func DisconnectMQTT() {
	if mqttClient != nil && mqttClient.IsConnected() {
		mqttClient.Disconnect(250)
		logger.Info("Disconnected from MQTT broker")
	}
}

// SendDiscoveryMessage sends Home Assistant discovery message
func SendDiscoveryMessage(cmd CommandConfig, clientID string) {
	deviceID := clientID
	sensorID := fmt.Sprintf("%s_%s", deviceID, sanitizeName(cmd.Name))

	discovery := HomeAssistantDiscovery{
		Name:       cmd.Name,
		StateTopic: fmt.Sprintf("homeassistant/sensor/%s/state", sensorID),
		UniqueID:   sensorID,
		Device: Device{
			Identifiers:  []string{deviceID},
			Name:         "Command Sensors",
			Model:        "HA Command to MQTT",
			Manufacturer: "Custom",
		},
	}

	if cmd.DeviceClass != "" {
		discovery.DeviceClass = cmd.DeviceClass
	}
	if cmd.Unit != "" {
		discovery.UnitOfMeasurement = cmd.Unit
	}
	if cmd.Icon != "" {
		discovery.Icon = cmd.Icon
	}
	if cmd.ForceUpdate {
		discovery.ForceUpdate = cmd.ForceUpdate
	}
	if cmd.StateClass != "" {
		discovery.StateClass = cmd.StateClass
	}
	if cmd.EntityCategory != "" {
		discovery.EntityCategory = cmd.EntityCategory
	}
	if cmd.ExpireAfter > 0 {
		discovery.ExpireAfter = cmd.ExpireAfter
	}

	payload, err := json.Marshal(discovery)
	if err != nil {
		logger.Errorf("Failed to marshal discovery message for %s: %v", cmd.Name, err)
		return
	}

	topic := fmt.Sprintf("homeassistant/sensor/%s/config", sensorID)
	token := mqttClient.Publish(topic, 0, true, payload)
	token.Wait()

	logger.Infof("Sent discovery message for %s", cmd.Name)
}

// PublishResult publishes command result to MQTT
func PublishResult(cmd CommandConfig, result string, clientID string) {
	deviceID := clientID
	sensorID := fmt.Sprintf("%s_%s", deviceID, sanitizeName(cmd.Name))
	topic := fmt.Sprintf("homeassistant/sensor/%s/state", sensorID)

	token := mqttClient.Publish(topic, 0, false, result)
	token.Wait()

	logger.Infof("Published result for %s: %s", cmd.Name, result)
}

// sanitizeName replaces spaces and special characters with underscores
func sanitizeName(name string) string {
	// Replace spaces and special characters with underscores
	result := strings.ToLower(name)
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ReplaceAll(result, "-", "_")
	// Remove any characters that aren't alphanumeric or underscore
	var sanitized strings.Builder
	for _, r := range result {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			sanitized.WriteRune(r)
		}
	}
	return sanitized.String()
}