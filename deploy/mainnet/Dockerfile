# Get stellar-core from the official image (as fallback)
FROM stellar/stellar-core:latest AS stellar-core

# Build environment - use golang with Ubuntu 20.04 (focal) base
FROM golang:1.21-focal AS build-env

# Install dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    curl \
    wget

# Final image - use Ubuntu 20.04 (focal) which is compatible with stellar-core
FROM golang:1.21-focal

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
    libldap-2.4-2 \
    libldap-common \
    libssl1.1 || \
    (echo "libssl1.1 not found in default repos, adding Debian repo" && \
     echo "deb http://security.debian.org/debian-security bullseye-security main" > /etc/apt/sources.list.d/debian-security.list && \
     apt-get update && apt-get install -y libssl1.1 libldap-2.4-2 libldap-common)

# Skip apt installation and use the stellar-core binary directly from the Docker image
RUN mkdir -p /usr/local/bin && \
    cp /tmp/stellar-core /usr/local/bin/stellar-core && \
    chmod +x /usr/local/bin/stellar-core

# Copy ALL libraries from the stellar-core image to ensure we have everything needed
COPY --from=stellar-core /usr/lib/x86_64-linux-gnu/*.so* /usr/lib/x86_64-linux-gnu/

# Copy specific libraries that might be missing
COPY --from=stellar-core /usr/bin/stellar-core /usr/local/bin/
RUN chmod +x /usr/local/bin/stellar-core

# Create symbolic links for libraries that might have different names
RUN if [ -f /usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2 ] && [ ! -f /usr/lib/x86_64-linux-gnu/libgssapi.so.3 ]; then \
        ln -s /usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2 /usr/lib/x86_64-linux-gnu/libgssapi.so.3; \
    fi

# Add Heimdal GSSAPI libraries manually
RUN apt-get update && apt-get install -y libgssapi3-heimdal || \
    (mkdir -p /tmp/heimdal && cd /tmp/heimdal && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libgssapi3-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     dpkg -x libgssapi3-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     cp -a usr/lib/x86_64-linux-gnu/libgssapi*.so* /usr/lib/x86_64-linux-gnu/ && \
     ldconfig)

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

# Check for LDAP libraries
RUN ls -la /usr/lib/x86_64-linux-gnu/libldap* || \
    (echo "libldap not found in expected location, downloading manually" && \
     mkdir -p /tmp/libldap && cd /tmp/libldap && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/o/openldap/libldap-2.4-2_2.4.49+dfsg-2ubuntu1.9_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/o/openldap/libldap-common_2.4.49+dfsg-2ubuntu1.9_all.deb && \
     dpkg -x libldap-2.4-2_2.4.49+dfsg-2ubuntu1.9_amd64.deb . && \
     dpkg -x libldap-common_2.4.49+dfsg-2ubuntu1.9_all.deb . && \
     cp -a usr/lib/x86_64-linux-gnu/libldap*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/liblber*.so* /usr/lib/x86_64-linux-gnu/ && \
     ldconfig && \
     ls -la /usr/lib/x86_64-linux-gnu/libldap*)

# Check for Heimdal GSSAPI libraries
RUN ls -la /usr/lib/x86_64-linux-gnu/libgssapi* || \
    (echo "libgssapi not found in expected location, downloading manually" && \
     mkdir -p /tmp/libgssapi && cd /tmp/libgssapi && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libgssapi3-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libheimntlm0-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libheimbase1-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libhx509-5-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libhcrypto4-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libroken18-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libasn1-8-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     wget http://nz2.archive.ubuntu.com/ubuntu/pool/main/h/heimdal/libwind0-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb && \
     dpkg -x libgssapi3-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     dpkg -x libheimntlm0-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     dpkg -x libheimbase1-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     dpkg -x libhx509-5-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     dpkg -x libhcrypto4-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     dpkg -x libroken18-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     dpkg -x libasn1-8-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     dpkg -x libwind0-heimdal_7.7.0+dfsg-1ubuntu1_amd64.deb . && \
     cp -a usr/lib/x86_64-linux-gnu/libgssapi*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libheimdal*.so* /usr/lib/x86_64-linux-gnu/ || true && \
     cp -a usr/lib/x86_64-linux-gnu/libhx509*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libhcrypto*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libroken*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libasn1*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libwind*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libheimbase*.so* /usr/lib/x86_64-linux-gnu/ && \
     cp -a usr/lib/x86_64-linux-gnu/libheimntlm*.so* /usr/lib/x86_64-linux-gnu/ && \
     ldconfig && \
     ls -la /usr/lib/x86_64-linux-gnu/libgssapi*)

# Install all possible dependencies that stellar-core might need
RUN apt-get update && apt-get install -y \
    libsasl2-2 \
    libsasl2-modules \
    libgnutls30 \
    libgssapi-krb5-2 \
    libk5crypto3 \
    libkrb5-3 \
    libkrb5support0 \
    libsasl2-modules-gssapi-mit \
    libheimntlm0-heimdal \
    libheimbase1-heimdal \
    libhx509-5-heimdal \
    libhcrypto4-heimdal \
    libroken18-heimdal \
    libasn1-8-heimdal \
    libwind0-heimdal \
    libgssapi3-heimdal || echo "Some dependencies could not be installed"

# Verify stellar-core binary works
RUN ls -la /usr/local/bin/stellar-core
RUN ldd /usr/local/bin/stellar-core || echo "Could not check stellar-core dependencies"

# Install any missing libraries identified by ldd
RUN ldd /usr/local/bin/stellar-core 2>&1 | grep "not found" | awk '{print $1}' > /tmp/missing_libs.txt || true
RUN if [ -s /tmp/missing_libs.txt ]; then \
        echo "Missing libraries detected:" && cat /tmp/missing_libs.txt; \
        for lib in $(cat /tmp/missing_libs.txt); do \
            apt-cache search $lib | grep -v dev; \
        done; \
    else \
        echo "No missing libraries detected"; \
    fi

# Install specific Heimdal libraries directly
RUN apt-get update && apt-get install -y \
    libgssapi3-heimdal \
    libheimntlm0-heimdal \
    libheimbase1-heimdal \
    libhx509-5-heimdal \
    libhcrypto4-heimdal \
    libroken18-heimdal \
    libasn1-8-heimdal \
    libwind0-heimdal || echo "Could not install some Heimdal libraries"

# Try to run stellar-core version
RUN /usr/local/bin/stellar-core version || echo "stellar-core version check failed"

# Set the entrypoint with direct path to stellar-core
ENTRYPOINT ["sh", "-c", "TZ=UTC/GMT /app/soroban-rpc --endpoint 0.0.0.0:80 --stellar-core-timeout \"50s\" --stellar-core-binary-path /usr/local/bin/stellar-core --captive-core-config-path /app/cmd/soroban-rpc/config/stellar-captive-core-live.toml --history-archive-urls https://history.stellar.org/prd/core-live/core_live_001,https://history.stellar.org/prd/core-live/core_live_002 --network-passphrase \"Public Global Stellar Network ; September 2015\""]
