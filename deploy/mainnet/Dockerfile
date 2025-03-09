# Build environment
FROM golang:1.21 AS build-env

# Install dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    curl \
    wget

# Final image
FROM golang:1.21

# Install required dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    wget \
    make \
    build-essential \
    git \
    libpq-dev \
    postgresql-client

# Remove potentially conflicting libc++ versions (if they exist)
RUN apt-get remove -y libc++1-10 libc++abi1-10 || true

# Add Stellar apt repository and install stellar-core
RUN wget -qO - https://apt.stellar.org/SDF.asc | apt-key add - \
    && echo "deb https://apt.stellar.org focal unstable" > /etc/apt/sources.list.d/SDF-unstable.list \
    && apt-get update \
    && apt-get install -y libc++1 libc++abi1

# Install Rust and Cargo
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="$PATH:/root/.cargo/bin"

# Create app directory
WORKDIR /app

# Install stellar-core from apt repository
RUN apt-get update && apt-get install -y stellar-core

# Ensure library paths are properly set
RUN ldconfig

# Verify stellar-core is working
RUN stellar-core version

# Copy the source code
COPY . .

# Set the working directory
WORKDIR /app

# Create the config directory
RUN mkdir -p /app/cmd/soroban-rpc/config

# Copy the captive core config
COPY deploy/mainnet/stellar-captive-core-live.toml /app/cmd/soroban-rpc/config/

# Build the libpreflight library
RUN cd /app/cmd/soroban-rpc/lib/preflight && cargo build --release

# Create a symbolic link to the library in a standard location
RUN ln -s /app/cmd/soroban-rpc/lib/preflight/target/release/libpreflight.so /usr/lib/libpreflight.so

# Set environment variables
ENV LD_LIBRARY_PATH="/app/cmd/soroban-rpc/lib/preflight/target/release:/usr/lib:/usr/lib/x86_64-linux-gnu"

# Set the entrypoint to run the binary directly
RUN cd /app && make build-soroban-rpc

# Set the entrypoint
ENTRYPOINT ["sh", "-c", "TZ=UTC/GMT /app/soroban-rpc --endpoint 0.0.0.0:80 --stellar-core-timeout \"50s\" --stellar-core-binary-path $(which stellar-core) --captive-core-config-path /app/cmd/soroban-rpc/config/stellar-captive-core-live.toml --history-archive-urls https://history.stellar.org/prd/core-live/core_live_001,https://history.stellar.org/prd/core-live/core_live_002 --network-passphrase \"Public Global Stellar Network ; September 2015\""]
