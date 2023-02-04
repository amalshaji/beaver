# Beaver

<img src="docs/beaver.png" height="250px">

## Client

Download the binary from releases or build from source. Refer `Makefile` for client build command. You may have to build for your architecture.

```bash
beaver - tunnel local ports to public URLs:

Usage:
      --config string      Config file path (default "/Users/amalshaji/.beaver/beaver_client.yaml")
      --subdomain string   Subdomain to tunnel http requests (default "<random_subdomain>")
      --port int           Local http server port (required)
```

### Config file for client

```yaml
targets :                 # Beaver server url (eg: wss://tunnel.example.com)
 - ws://127.0.0.1:8080
poolidlesize : 1          # Default number of concurrent open (TCP) connections to keep idle per WSP server(optional)
poolmaxsize : 100         # Maximum number of concurrent open (TCP) connections per WSP server(optional)
secretkey : ThisIsASecret # Users secret key set in the server config
```

## Server

```bash
docker pull amalshaji/beaver:latest
```

### Config file for server

```yaml
host : 127.0.0.1             # Address to bind the HTTP server
port : 8080                  # Port to bind the HTTP server
timeout : 3000               # Time to wait before acquiring a WS connection to forward the request (milliseconds, optional)
idletimeout : 60000          # Time to wait before closing idle connection when there is enough idle connections (milliseconds, optional)
users:                       # User specific secret keys
  - identifier: foo@bar.com
    secretkey: ThisIsASecret
  - identifier: max@xam.com
    secretkey: ThisIsASecret@2
```

## Credits

This project is a fork of [hgsgtk/wsp](https://github.com/hgsgtk/wsp)
