# WG Controller

A wireguard VPN server and control plane with central web management.

## Features

- Easily host your own Wireguard server with Docker or Kubernetes
- Manage users and devices from a modern web interface
- Integrated DNS server resolves devices by their configured name
- Internal IP routing between clients
- Share access to remote networks on a client by client basis
- Synchronization of keys and settings between clients and server (using [wg-controller-client](https://github.com/wg-controller/wg-controller-client))
- Support for standard wireguard clients and 3rd party devices

## Screenshots

![Clients Page](/screenshot1.png?raw=true "Client Management Page")

## Deployment

### Docker

- Clone repo or download docker-compose.yaml
- Generate WG_PRIVATE_KEY and DB_AES_KEY to fill environment fields in docker-compose.yaml

  ```
  docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-wg-key
  ```

  ```
  docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-db-key
  ```

- Fill in remaining environment fields in docker-compose.yaml
- Start server with docker compose

  ```
  docker compose up
  ```

- The web interface will be running on port :8080

> [!WARNING]
> Do not host this on the internet without an appropriate SSL reverse proxy (see [NGINX](https://hub.docker.com/_/nginx), [Caddy](https://caddyserver.com) etc)

### Kubernetes

- Clone repo or download kube-manifests.yaml
- Generate WG_PRIVATE_KEY and DB_AES_KEY to fill env fields in kube-manifests.yaml

  ```
  docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-wg-key
  ```

  ```
  docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-db-key
  ```

- Fill in remaining env fields in kube-manifests.yaml
- Configure ingress domain, SSL etc
- Deploy to kubernetes with kubectl

  ```
  kubectl apply -f kube-manifests.yaml
  ```

## Development

[Tygo](https://github.com/gzuidhof/tygo) is used for generating TypeScript types from Golang types <br>
Install Tygo with `go install github.com/gzuidhof/tygo@latest` <br>
Running `tygo generate` will export Go types to frontend.
