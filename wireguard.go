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

	"github.com/lampy255/net-tbm/db"
	"github.com/lampy255/net-tbm/types"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var wireguard_cmd *exec.Cmd
var wg *wgctrl.Client

// Gets wireguard-go version information
func GetWireguard() (version string, err error) {
	cmd := exec.Command("wireguard-go", "--version")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Starts wireguard-go with a goroutine attached
func StartWireguard() {
	cmd := exec.Command("wireguard-go", "-f", ENV.INTERFACE_NAME)
	wireguard_cmd = cmd

	go func() {
		log.Println("Starting wireguard-go")
		op, err := cmd.Output()
		if err != nil {
			fmt.Println(string(op))
			os.Exit(1)
		}
	}()

	// Init wg-ctrl
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
