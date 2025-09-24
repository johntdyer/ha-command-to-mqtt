ARG BUILD_FROM
FROM $BUILD_FROM

# Install necessary packages
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    openssh-client \
    bash \
    curl

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Copy the pre-built binary from the host
COPY ha-command-to-mqtt /usr/local/bin/

# Make it executable
RUN chmod +x /usr/local/bin/ha-command-to-mqtt

# Copy run script
COPY run.sh /
RUN chmod a+x /run.sh

# Create directories for configuration and SSH keys
RUN mkdir -p /config /share/ssh_keys

# Labels
LABEL \
    io.hass.name="Command to MQTT" \
    io.hass.description="Execute commands and publish results to MQTT for Home Assistant" \
    io.hass.arch="${BUILD_ARCH}" \
    io.hass.type="addon" \
    io.hass.version="${BUILD_VERSION}" \
    maintainer="jdyer" \
    org.opencontainers.image.title="Command to MQTT" \
    org.opencontainers.image.description="Execute commands and publish results to MQTT for Home Assistant" \
    org.opencontainers.image.vendor="jdyer" \
    org.opencontainers.image.authors="jdyer" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.source="https://github.com/jdyer/ha-command-to-mqtt"

CMD ["/run.sh"]