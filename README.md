# Beaver



<p align="center">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/amalshaji/beaver/unit-tests.yml">
    <img alt="GitHub" src="https://img.shields.io/github/license/amalshaji/beaver">
    <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/amalshaji/beaver">
</p>

<p align="center">
    <img src="docs/beaver.png" height="250px">
</p>

> **Warning**
> This project is in a very early stage, If you find any bugs, please raise an issue.

## Client

### Install

Using homebrew

```shell
➜ brew tap amalshaji/taps
➜ brew install beaver
```

Or, download the binary from the [releases page](https://github.com/amalshaji/beaver/releases)

### Usage

Once the client is installed, run:

```shell
➜ beaver config --init
```

This creates a basic config at `$HOME/.beaver/beaver_client.yaml`,

```yaml
target: 
secretkey: 
tunnels:
  - name: tunnel-1
    subdomain: subdomain-1
    port: 8000
```

Update your `target` and `secretKey`, and you're ready to go.

## Server

> [Deploying the server using caddy and cloudflare](https://github.com/amalshaji/beaver/wiki/Deploying-the-server-using-caddy)

### Deploy

1. Using docker

    ```shell
    docker run \
      -v $PWD/docs/beaver_server.yaml:/app/config/beaver_server.yaml \
      -v $PWD/data:/app/data/ \
      -p 8080:8080 --restart unless-stopped amalshaji/beaver:0.1.0
    ```

    Replace `$PWD/docs/beaver_server.yaml` with path to your config file

1. Using docker compose

    ```yaml
    services:
      beaver:
        image: amalshaji/beaver:0.1.0
        volumes:
          - ./docs/beaver_server.yaml:/app/config/beaver_server.yaml
          - ./data:/app/data/
        ports:
          - 8080:8080
        restart: unless-stopped
    ```

    Start the server:

    ```shell
    docker compose up -d
    ```

### Config

```yaml
# beaver_server.yaml

host : 0.0.0.0                  # Address to bind the HTTP server
port : 8080                     # Port to bind the HTTP server
domain: localhost               # Domain on which the server will be running (eg: tunnel.example.com)            
secure: false                   # Whether the server runs under https
timeout : 3000                  # Time to wait before acquiring a WS connection to forward the request (milliseconds)
idletimeout : 60000             # Time to wait before closing idle connection when there is enough idle connections (milliseconds)
```

## Credits

This project is a fork of [hgsgtk/wsp](https://github.com/hgsgtk/wsp)
