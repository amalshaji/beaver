# Beaver

<img src="docs/beaver.png" height="250px">

## Client

Download the binary from [releases page](https://github.com/amalshaji/beaver/releases).

```bash
➜  beaver --help
Tunnel local ports to public URLs

Usage:
  beaver [flags]
  beaver [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  http        Tunnel local http servers
  start       Start tunnels defined in the config file

Flags:
      --config string   Path to the client config file (default "/Users/amalshaji/.beaver/beaver_client.yaml")
  -h, --help            help for beaver

Use "beaver [command] --help" for more information about a command.

---

➜  beaver http --help
Tunnel local http servers

Usage:
  beaver http [PORT] [flags]

Flags:
  -h, --help               help for http
      --subdomain string   Subdomain to tunnel http requests (default "<random_subdomain>")

---

➜  beaver start --help
Start tunnels defined in the config file

Usage:
  beaver start [--all] or [tunnel1 tunnel2] [flags]

Flags:
      --all    Start all tunnels listed in the config
  -h, --help   help for start
```

#### Example

```shell
➜  beaver git:(main) ✗ ./beaver --config docs/beaver_client.yaml --port 8000
2023/02/05 19:46:07 Creating tunnel connection
2023/02/05 19:46:07 Tunnel running on https://sccrej.tunnel.example.com
```

Now, `https://sccrej.tunnel.example.com -> http://localhost:8000`

#### Config

```yaml
targets:                 
  - ws://127.0.0.1:8080    # Beaver server url (eg: wss://tunnel.example.com)
poolidlesize: 1            # Default number of concurrent open (TCP) connections to keep idle per WSP server(optional)
poolmaxsize: 100           # Maximum number of concurrent open (TCP) connections per WSP server(optional)
secretkey: ThisIsASecret   # User's secret key set in the server config
```

## Server

Download the binary from [releases page](https://github.com/amalshaji/beaver/releases), or use the [docker image](https://hub.docker.com/r/amalshaji/beaver)

#### Deploy

1. Using the binary

    ```shell
    ./beaver_server --config docs/beaver_server.yaml
    ```

1. Using docker

    ```shell
    docker run \
      -v $PWD/docs/beaver_server.yaml:/app/config/beaver_server.yaml \
      -p 8080:8080 amalshaji/beaver:latest
    ```

    Replace `$PWD/docs/beaver_server.yaml` with path to your config file

1. Using docker compose

    ```yaml
    services:
      beaver:
        image: amalshaji/beaver:latest
        volumes:
          - ./docs/beaver_server.yaml:/app/config/beaver_server.yaml
        ports:
          - 8080:8080
        restart: unless-stopped
    ```

    Start the server:

    ```shell
    docker compose up -d
    ```

#### Config

```yaml
host : 0.0.0.0                  # Address to bind the HTTP server
port : 8080                     # Port to bind the HTTP server
domain: localhost               # Domain on which the server will be running (eg: tunnel.example.com)            
secure: false                   # Whether the server runs under https
timeout : 3000                  # Time to wait before acquiring a WS connection to forward the request (milliseconds)
idletimeout : 60000             # Time to wait before closing idle connection when there is enough idle connections (milliseconds)
users:
  - identifier: foo@bar.com
    secretkey: ThisIsASecret
  - identifier: max@xam.com
    secretkey: ThisIsASecret@2

```

## Credits

This project is a fork of [hgsgtk/wsp](https://github.com/hgsgtk/wsp)

## Checkout [wiki](https://github.com/amalshaji/beaver/wiki) for examples
