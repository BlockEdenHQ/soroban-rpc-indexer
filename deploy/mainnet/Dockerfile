# Get stellar-core from the official image (as fallback)
FROM stellar/stellar-core:latest AS stellar-core

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

# Install Rust and Cargo
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="$PATH:/root/.cargo/bin"

# Create app directory
WORKDIR /app

# Copy stellar-core binary from the stellar-core image as fallback
COPY --from=stellar-core /usr/bin/stellar-core /tmp/stellar-core

# Install libc++ dependencies manually (to avoid version conflicts)
RUN apt-get update && apt-get install -y \
    libc++1 \
    libc++abi1 \
    libboost-all-dev \
    libsodium-dev \
    libsqlite3-dev \
    libssl1.1 || \
    (echo "libssl1.1 not found in default repos, adding Debian repo" && \
     echo "deb http://security.debian.org/debian-security bullseye-security main" > /etc/apt/sources.list.d/debian-security.list && \
     apt-get update && apt-get install -y libssl1.1)

# Add Stellar apt repositories (both stable and unstable)
RUN wget -qO - https://apt.stellar.org/SDF.asc | apt-key add - \
    && echo "deb https://apt.stellar.org focal stable" > /etc/apt/sources.list.d/SDF-stable.list \
    && echo "deb https://apt.stellar.org focal unstable" > /etc/apt/sources.list.d/SDF-unstable.list

# Update apt and try to install stellar-core with specific version if available
RUN apt-get update && apt-get install -y stellar-core || \
    (apt-get update && apt-cache madison stellar-core && \
     apt-get install -y stellar-core=20.2.0-1716.rc3.34d82fc00.focal || \
     apt-get install -y stellar-core || \
     (echo "Falling back to copying stellar-core from Docker image" && \
      mkdir -p /usr/local/bin && \
      cp /tmp/stellar-core /usr/local/bin/stellar-core && \
      chmod +x /usr/local/bin/stellar-core))

# Copy necessary libraries from the stellar-core image as fallback
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/libc++*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/libpq*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/libboost*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/libsodium*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/libsqlite*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/libssl*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/libcrypto*.so* /usr/lib/x86_64-linux-gnu/

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

# Check for OpenSSL library
RUN ls -la /usr/lib/x86_64-linux-gnu/libssl* || \
    (echo "libssl not found in expected location, downloading manually" && \
     mkdir -p /tmp/libssl && cd /tmp/libssl && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.1f-1ubuntu2.19_amd64.deb && \
     dpkg -x libssl1.1_1.1.1f-1ubuntu2.19_amd64.deb . && \
     cp -a usr/lib/x86_64-linux-gnu/libssl.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libcrypto.so* /usr/lib/x86_64-linux-gnu/ && \
     ldconfig && \
     ls -la /usr/lib/x86_64-linux-gnu/libssl*)

# Find stellar-core binary and verify it works
RUN which stellar-core || which /usr/local/bin/stellar-core || echo "stellar-core not found in PATH"
RUN (stellar-core version || /usr/local/bin/stellar-core version) || echo "stellar-core version check failed"

# Set the entrypoint with fallback paths for stellar-core
ENTRYPOINT ["sh", "-c", "CORE_PATH=$(which stellar-core 2>/dev/null || echo /usr/local/bin/stellar-core) && \
TZ=UTC/GMT /app/soroban-rpc --endpoint 0.0.0.0:80 --stellar-core-timeout \"50s\" --stellar-core-binary-path $CORE_PATH --captive-core-config-path /app/cmd/soroban-rpc/config/stellar-captive-core-live.toml --history-archive-urls https://history.stellar.org/prd/core-live/core_live_001,https://history.stellar.org/prd/core-live/core_live_002 --network-passphrase \"Public Global Stellar Network ; September 2015\""]
