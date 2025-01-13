package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/wg-controller/wg-controller/db"
	"github.com/wg-controller/wg-controller/types"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var wireguard_cmd *exec.Cmd
var wg *wgctrl.Client

// Starts wireguard-go with a goroutine attached
func StartWireguard() {
	cmd := exec.Command("wireguard-go", "-f", ENV.WG_INTERFACE)
	wireguard_cmd = cmd

	go func() {
		log.Println("Starting wireguard-go")
		op, err := cmd.Output()
		if err != nil {
			fmt.Println(string(op))
			os.Exit(1)
		}
	}()

	// Create wireguard client
	client, err := wgctrl.New()
	if err != nil {
		log.Fatal("Unable to connect to wireguard-go")
	}
	wg = client
}

// Kills wireguard-go
func StopWireguard() {
	if wireguard_cmd != nil {
		err := wireguard_cmd.Process.Kill()
		if err != nil {
			log.Fatal("Error stopping wireguard-go:", err)
		}
	}
}

func PruneWireguardPeers(peerPublicKeys []string) error {
	var wgPeers []wgtypes.PeerConfig
	for _, pk := range peerPublicKeys {
		// Parse PublicKey
		publicKey, err := wgtypes.ParseKey(pk)
		if err != nil {
			return err
		}

		// Create wireguard-go peer configuration with Remove set to true
		wgPeer := wgtypes.PeerConfig{
			PublicKey: publicKey,
			Remove:    true,
		}
		wgPeers = append(wgPeers, wgPeer)
	}

	// Append the wireguard-go configuration
	return wg.ConfigureDevice(ENV.WG_INTERFACE, wgtypes.Config{
		ReplacePeers: false,
		Peers:        wgPeers,
	})
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
		// Skip peers that are not enabled
		if !peer.Enabled {
			continue
		}

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
		// Append peer's own subnet
		_, ipNet, err := net.ParseCIDR(peer.RemoteTunAddress + "/32")
		if err != nil {
			log.Println("Error parsing peer's own subnet:", err)
		} else {
			allowedIPs = append(allowedIPs, *ipNet)
		}

		// Create wireguard-go peer configuration
		wgPeer := wgtypes.PeerConfig{
			PublicKey:                   publicKey,
			PresharedKey:                &preSharedKey,
			PersistentKeepaliveInterval: &keepAliveDuration,
			AllowedIPs:                  allowedIPs,
			ReplaceAllowedIPs:           true,
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

	// Append the wireguard-go configuration
	return wg.ConfigureDevice(ENV.WG_INTERFACE, wgtypes.Config{
		ReplacePeers: false,
		Peers:        wgPeers,
		PrivateKey:   &privateKey,
		ListenPort:   &wgPort,
	})
}

func GetWireguardPeer(storedPeer types.Peer) (types.Peer, error) {
	// Get wireguard data
	device, err := wg.Device(ENV.WG_INTERFACE)
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

	return types.Peer{}, errors.New("peer not found in wireguard-go")
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
