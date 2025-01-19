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
```
docker run -e PUBLIC_HOST='wg.example.com' \
-e ADMIN_EMAIL='admin@example.com' \
-e ADMIN_PASS='examplepass1234!!' \
-e WG_PRIVATE_KEY='' \
-e DB_AES_KEY='' \
-p 8081:8081 \
-p 51820:51820/udp \
--cap-add NET_ADMIN \
--cap-add SYS_MODULE \
--sysctl 'net.ipv4.conf.all.src_valid_mark=1' \
--sysctl 'net.ipv4.ip_forward=1' \
--name wg-controller ghcr.io/wg-controller/wg-controller:latest
```

### Kubernetes
```
kubectl apply -f kube-manifests.yaml
```



## Development
[Tygo](https://github.com/gzuidhof/tygo) is used for generating TypeScript types from Golang types <br>
Install Tygo with `go install github.com/gzuidhof/tygo@latest` <br>
Running `tygo generate` will export Go types to frontend.
