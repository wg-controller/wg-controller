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

## Deployment

### Docker

Use following commands to generate WG_PRIVATE_KEY and DB_AES_KEY in docker-compose.yaml

```
docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-wg-key
```

```
docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-db-key
```

Start server with docker compose

```
docker compose up
```

### Kubernetes

Use following commands to generate WG_PRIVATE_KEY and DB_AES_KEY in kube-manifests.yaml

```
docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-wg-key
```

```
docker run --rm -it ghcr.io/wg-controller/wg-controller:latest /app/main generate-db-key
```

Deploy to kubernetes with kubectl

```
kubectl apply -f kube-manifests.yaml
```

## Development

[Tygo](https://github.com/gzuidhof/tygo) is used for generating TypeScript types from Golang types <br>
Install Tygo with `go install github.com/gzuidhof/tygo@latest` <br>
Running `tygo generate` will export Go types to frontend.
