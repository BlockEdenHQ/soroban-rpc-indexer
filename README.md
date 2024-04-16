# soroban-rpc
RPC Server for Soroban Contracts.

# soroban-indexer

This project is for the indexing service of BlockEden x Stellar.
- [BlockEden x Stellar](https://docs.google.com/document/d/1hXvILdI2SKhgR3cv_xVVdVevI1U4dEsaYgzeFy0XwpU/edit?pli=1): PRFAQ BlockEden indexer for Stellar.


# setup

```bash
wget -qO - https://apt.stellar.org/SDF.asc | sudo apt-key add -
echo "deb https://apt.stellar.org $(lsb_release -cs) stable" | sudo tee -a /etc/apt/sources.list.d/SDF.list
sudo apt update
sudo apt install stellar-core
```

Or if you want to build from source, checkout `stellar-core` to branch `v20.3.0`，then build and install：

prepare the environment https://github.com/stellar/stellar-core/blob/master/INSTALL.md#ubuntu and then

```bash
cd ~
git clone https://github.com/stellar/stellar-core.git
cd stellar-core

git checkout TODO_version

git submodule init
git submodule update
./autogen.sh
./configure
make -j6 
make install
```



2. build and run this indexer



```bash
cd ~
git clone git@github.com:BlockEdenHQ/soroban-indexer.git
cd soroban-indexer
make build

# prepare database connection URL
echo 'POSTGRES_DSN="TODO"' > .env
make migrate # create database tables

# run
make dev-ubuntu-mainnet

# run in background
sudo nohup make dev-ubuntu-mainnet < /dev/null > output.log 2>&1 &
disown

# or infinitely run
nohup ./run.sh &
```

How to connect to PostgreSQL in gcp? create instance, create database, create user, allow certain IP.

## How to integrate with testnet and livenet?

This project is for futurenet use mainly, but can also be configured to use in testnet and mainnet. Note that there are no events data on testnet and mainnet. Please check Makefile commands `dev-ubuntu-*` for details.

## How to install Golang & Rust?

Install build tools first https://github.com/stellar/stellar-core/blob/master/INSTALL.md#ubuntu

```
# golang
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go

# rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```



## How to run sorobon-indexer on Google cloud VM instance?

### Expose 80

1. install stellar-core binary described as the above to the same directory with `sorobon-indexer`.
2. set up ubuntu firewall `sudo ufw disable `
3. run `dev-ubuntu` to listen on `:80` port
4. in Google Cloud terminal, mount the network and firewall to allow `:80`

Then visit the public IP with `:80` and see if it works.


## How to run in production

> 🔥 remember to rebuild after pulling the latest change


```
make build-soroban-rpc
```


```bash
sudo nohup make dev-ubuntu-testnet < /dev/null > output.log 2>&1 &
disown
```



















# Soroban-RPC

Soroban-RPC allows you to communicate directly with Soroban via a JSON RPC interface.

For example, you can build an application and have it send a transaction, get ledger and event data or simulate transactions.

## Dependencies
- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/doc/install)
- [Rust](https://www.rust-lang.org/tools/install)
- [Cargo](https://doc.rust-lang.org/cargo/getting-started/installation.html)

## Building Stellar-Core
Soroban-RPC requires an instance of stellar-core binary on the same host. This is referred to as the `Captive Core`.
Since, we are building RPC from source, we recommend considering two approaches to get the stellar-core binary:
- If saving time is top priority and your development machine is on a linux debian OS, then consider installing the
  testnet release candidates from the [testing repository.](https://apt.stellar.org/pool/unstable/s/stellar-core/)
- The recommended option is to compile the core source directly on your machine:
    - Clone the stellar-core repo:
        ```bash
        git clone https://github.com/stellar/stellar-core.git
        cd stellar-core
        ```
    - Fetch the tags and checkout the testnet release tag:
        ```bash
        git fetch --tags
        git checkout tags/v20.0.0-rc.2.1 -b soroban-testnet-release
        ```
    - Follow the build steps listed in [INSTALL.md](https://github.com/stellar/stellar-core/blob/master/INSTALL.md) file for the instructions on building the local binary

## Building Soroban-RPC
- Similar to stellar-core, we will clone the soroban-tools repo and checkout the testnet release tag:
```bash
git clone https://github.com/stellar/soroban-tools.git
cd soroban-tools
git fetch --tags
git checkout tags/v20.0.0-rc4 -b soroban-testnet-release
```
- Build soroban-rpc target:
```bash
make build-soroban-rpc
```
This will install and build the required dependencies and generate a `soroban-rpc` binary in the working directory.

## Configuring and Running RPC Server
- Both stellar-core and soroban-rpc require configuration files to run.
    - For production, we specifically recommend running Soroban RPC with a TOML configuration file rather than CLI flags.
    - There is a new subcommand `gen-config-file` which takes all the same arguments as the root command (or no arguments at all),
      and outputs the resulting config toml file to stdout.
        ```bash
        ./soroban-rpc gen-config-file
        ```
    - Paste the output to a file and save it as `.toml` file in any directory.
    - Make sure to update the config values to testnet specific ones. You can refer to [Configuring](https://docs.google.com/document/d/1SIbrFWFgju5RAsi6stDyEtgTa78VEt8f3HhqCLoySx4/edit#heading=h.80d1jdtd7ktj) section in the Runbook for specific config settings.
- If everything is set up correctly, then you can run the RPC server with the following command:
```bash
./soroban-rpc --config-path <PATH_TO_THE_RPC_CONFIG_FILE>
```
