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
> Do not host this on the internet without an appropriate SSL reverse proxy (see [NGINX](https://hub.docker.com/_/nginx), [Caddy](https://caddyserver.com))

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

## Options

| Env              | Default       | Example                                      |
| ---------------- | ------------- | -------------------------------------------- |
| PUBLIC_HOST      | required      | wg.example.com                               |
| ADMIN_EMAIL      | required      | admin@example.com                            |
| ADMIN_PASS       | required      | SuP3Rs8cureP4ssw0rd#                         |
| WG_PRIVATE_KEY   | required      | WFgLw2vV1Pc1EhtRXdFNHOopmuNl9GZluRFhI73Mf2o= |
| DB_AES_KEY       | required      | CQLZLLfq+XXQKWrLDDvy0vine6Yil3SGxGJEUHK32yU= |
| SERVER_CIDR      | 172.19.0.0/24 | 192.168.10.0/24                              |
| SERVER_ADDRESS   | 172.19.0.254  | 192.168.10.1                                 |
| EGRESS_INTERFACE | eth0          | eth2                                         |
| WG_INTERFACE     | wg0           | UTUN11                                       |
| WG_PORT          | 51820         | 51821                                        |
| API_PORT         | 8081          | 9000                                         |
| SERVER_HOSTNAME  | wg-controller | my-vpn-server                                |
| UPSTREAM_DNS     | 8.8.8.8       | 1.1.1.1                                      |

## Security

- Wireguard keys encrypted at rest with AES256
- Passwords and API keys salted and hashed before storage

## Development

[Tygo](https://github.com/gzuidhof/tygo) is used for generating TypeScript types from Golang types <br>
Install Tygo with `go install github.com/gzuidhof/tygo@latest` <br>
Running `tygo generate` will export Go types to frontend.
