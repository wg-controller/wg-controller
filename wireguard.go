package main

import (
	"errors"
	"net"
	"os/exec"
	"strconv"
	"time"

	"github.com/lampy255/wg-controller/db"
	"github.com/lampy255/wg-controller/types"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var wg *wgctrl.Client

func InitWireguardInterface() error {
	// Create wireguard client
	client, err := wgctrl.New()
	if err != nil {
		return err
	}
	wg = client

	// Create wireguard interface
	cmd := exec.Command("wg-quick", "up", ENV.INTERFACE_NAME)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func SyncWireguardConfiguration() error {
	// Get all peers from the database
	peers, err := db.GetPeers()
	if err != nil {
		return err
	}

	// Convert peers to wireguard-go peer configurations
	var wgPeers []wgtypes.PeerConfig
	for _, peer := range peers {
		// Convert KeepAliveSeconds to time.Duration
		keepAliveDuration := time.Duration(peer.KeepAliveSeconds) * time.Second

		// Parse PublicKey
		publicKey, err := wgtypes.ParseKey(peer.PublicKey)
		if err != nil {
			return err
		}

		// Parse PreSharedKey
		preSharedKey, err := wgtypes.ParseKey(peer.PreSharedKey)
		if err != nil {
			return err
		}

		// Parse allowed subnets
		allowedIPs := []net.IPNet{}
		for _, subnet := range peer.RemoteSubnets {
			_, ipNet, err := net.ParseCIDR(subnet)
			if err != nil {
				break
			}
			allowedIPs = append(allowedIPs, *ipNet)
		}

		// Create wireguard-go peer configuration
		wgPeer := wgtypes.PeerConfig{
			PublicKey:                   publicKey,
			PresharedKey:                &preSharedKey,
			PersistentKeepaliveInterval: &keepAliveDuration,
			AllowedIPs:                  allowedIPs,
		}
		wgPeers = append(wgPeers, wgPeer)
	}

	// Parse ENV.WG_PRIVATE_KEY
	privateKey, err := wgtypes.ParseKey(ENV.WG_PRIVATE_KEY)
	if err != nil {
		return err
	}

	// Convert ENV.WG_PORT to int
	wgPort, err := strconv.Atoi(ENV.WG_PORT)
	if err != nil {
		return err
	}

	// Overwrite the wireguard-go configuration
	return wg.ConfigureDevice(ENV.INTERFACE_NAME, wgtypes.Config{
		ReplacePeers: true,
		Peers:        wgPeers,
		PrivateKey:   &privateKey,
		ListenPort:   &wgPort,
	})
}

func GetWireguardPeer(storedPeer types.Peer) (types.Peer, error) {
	// Get wireguard data
	device, err := wg.Device(ENV.INTERFACE_NAME)
	if err != nil {
		return types.Peer{}, err
	}

	// Find stored peer in wireguard peers
	for _, wgPeer := range device.Peers {
		if wgPeer.PresharedKey.String() == storedPeer.PreSharedKey {
			// Create extended peer
			storedPeer.TransmitBytes = wgPeer.TransmitBytes
			storedPeer.ReceiveBytes = wgPeer.ReceiveBytes
			storedPeer.LastSeenUnixMillis = wgPeer.LastHandshakeTime.UnixMilli()
			if wgPeer.Endpoint != nil {
				storedPeer.LastIPAddress = wgPeer.Endpoint.IP.String()
			}
			return storedPeer, nil
		}
	}

	return types.Peer{}, errors.New("peer not found")
}

func NewWireguardPrivateKey() (privKey string, err error) {
	// Generate new key pair
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", err
	}

	return key.String(), nil
}

func NewWireguardPreSharedKey() (preSharedKey string, err error) {
	key, err := wgtypes.GenerateKey()
	if err != nil {
		return "", err
	}

	return key.String(), nil
}

func GetWireguardPublicKey(privateKey string) (pubKey string, err error) {
	key, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return "", err
	}

	return key.PublicKey().String(), nil
}
