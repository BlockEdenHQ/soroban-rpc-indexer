# Soroban RPC Mainnet Docker Setup

This directory contains Docker Compose configuration for running Soroban RPC service for the Stellar mainnet.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

### 1. Build and Start the Service

```bash
# Navigate to this directory
cd /path/to/soroban-rpc-indexer/deploy/mainnet

# Build and start the service
docker compose up -d --build
```

### 2. Check the Status

```bash
# Check if the container is running
deploy-compose ps

# View logs
deploy-compose logs -f
```

### 3. Stop the Service

```bash
deploy-compose down
```

## Configuration

- The Soroban RPC service is exposed on port 80
- The service uses the stellar-core binary in 'captive core' mode with the configuration in `stellar-captive-core-live.toml`
- Data is persisted in Docker volumes

## Notes

- This setup builds a custom image for Soroban RPC that includes the stellar-core binary
- The stellar-core binary is used by soroban-rpc in 'captive core' mode, but no separate stellar-core service is needed
- The configuration is based on the `dev-ubuntu-mainnet` target in the project's Makefile
