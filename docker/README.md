# Docker Setup for HA Command to MQTT

This directory contains the Docker configuration files for running HA Command to MQTT in containers.

## Files

- `Dockerfile` - Multi-stage Docker build for the Go application
- `docker-compose.yml` - Complete setup with MQTT broker and application

## Quick Start

1. **Build and run with docker-compose:**

   ```bash
   cd docker
   docker-compose up -d
   ```

2. **Build Docker image manually:**

   ```bash
   cd docker
   docker build -f Dockerfile -t ha-command-to-mqtt ..
   ```

3. **Run container with custom config:**

   ```bash
   docker run -v /path/to/config.yaml:/root/config.yaml:ro ha-command-to-mqtt
   ```

## Configuration

The docker-compose setup includes:

- **ha-command-to-mqtt**: The main application container
- **mosquitto**: MQTT broker for local testing
- **Volumes**: Persistent storage for MQTT broker data
- **Network**: Shared network for container communication

### Environment Variables

You can configure the application using environment variables instead of a config file. See the commented examples in `docker-compose.yml`.

### Config File Mount

By default, the compose file mounts `../config.yaml` into the container. Make sure you have a valid config file in the project root, or update the volume mount path.

## Production Considerations

- Update the `networks.homeassistant.external: true` to match your actual network setup
- Consider using Docker secrets for sensitive information like MQTT passwords
- Review and customize the mosquitto configuration file (`../mosquitto.conf`)
- Set appropriate resource limits for production deployments
