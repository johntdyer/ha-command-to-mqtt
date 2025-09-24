#!/bin/bash

# Home Assistant Add-on Build Script
# This script builds the Go application for multiple architectures
# and prepares it for Home Assistant add-on deployment

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADDON_DIR="${PROJECT_DIR}/addon"

echo "Building Command to MQTT add-on..."
echo "Project directory: ${PROJECT_DIR}"
echo "Add-on directory: ${ADDON_DIR}"

# Supported architectures for Home Assistant add-ons
ARCHITECTURES=(
    "linux-amd64"
    "linux-arm64"
    "linux-armv7"
    "linux-armv6"
    "linux-386"
)

# Clean previous builds
echo "Cleaning previous builds..."
rm -f "${ADDON_DIR}/ha-command-to-mqtt"*

# Build for local architecture first (for testing)
echo "Building for local testing..."
cd "${PROJECT_DIR}"
go build -o "${ADDON_DIR}/ha-command-to-mqtt" -ldflags="-w -s" .

echo "✅ Local build complete: ${ADDON_DIR}/ha-command-to-mqtt"

# Function to build for specific architecture
build_arch() {
    local arch=$1
    local output_name="ha-command-to-mqtt-${arch}"

    echo "Building for ${arch}..."
    cd "${PROJECT_DIR}"

    # Set Go environment based on architecture
    case $arch in
        "linux-amd64")
            env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
                -o "${ADDON_DIR}/${output_name}" -ldflags="-w -s" .
            ;;
        "linux-arm64")
            env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
                -o "${ADDON_DIR}/${output_name}" -ldflags="-w -s" .
            ;;
        "linux-armv7")
            env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build \
                -o "${ADDON_DIR}/${output_name}" -ldflags="-w -s" .
            ;;
        "linux-armv6")
            env GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build \
                -o "${ADDON_DIR}/${output_name}" -ldflags="-w -s" .
            ;;
        "linux-386")
            env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build \
                -o "${ADDON_DIR}/${output_name}" -ldflags="-w -s" .
            ;;
        *)
            echo "❌ Unknown architecture: ${arch}"
            return 1
            ;;
    esac

    echo "✅ Built: ${output_name}"
}

# Build for all architectures
if [ "$1" = "--all-arch" ]; then
    echo "Building for all Home Assistant architectures..."
    for arch in "${ARCHITECTURES[@]}"; do
        build_arch "$arch"
    done
    echo ""
    echo "All architecture builds complete!"
    echo "Files in ${ADDON_DIR}:"
    ls -la "${ADDON_DIR}"/ha-command-to-mqtt*
fi

echo ""
echo "Add-on files ready in: ${ADDON_DIR}"
echo ""
echo "Next steps:"
echo "1. Copy the addon/ directory to your Home Assistant add-ons repository"
echo "2. Install and configure the add-on through Home Assistant"
echo "3. For multi-architecture support, use Docker buildx with the provided binaries"

# Create a simple repository.yaml for local add-on development
if [ ! -f "${ADDON_DIR}/../repository.yaml" ]; then
    cat > "${ADDON_DIR}/../repository.yaml" << EOF
name: "Local Add-ons Repository"
url: "https://github.com/your-username/ha-addons"
maintainer: "Your Name <your.email@example.com>"
EOF
    echo ""
    echo "📝 Created repository.yaml for local development"
fi

echo ""
echo "🏠 Home Assistant Add-on build complete!"