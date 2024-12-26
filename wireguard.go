package main

import (
	"log"
	"os/exec"

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
		err := cmd.Run()
		if err != nil {
			log.Fatal("Error starting wireguard-go:", err)
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

// Overwrites wireguard-go configuration with supplied peers
func OverwritePeers(peers []types.Peer) error {
	wg.ConfigureDevice(ENV.INTERFACE_NAME, wgtypes.Config{
		ReplacePeers: true,
	})

	return nil
}

func GeneratePeerConfig(peer types.Peer) (wgtypes.PeerConfig, error) {
	return wgtypes.PeerConfig{}, nil
}

func GetWireguardPeers(storedPeers []types.Peer) ([]types.PeerExtended, error) {
	// Get wireguard data
	device, err := wg.Device(ENV.INTERFACE_NAME)
	if err != nil {
		return nil, err
	}

	// Loop through stored peers and match with wireguard peers
	var ExtendedPeers []types.PeerExtended
	for _, storedPeer := range storedPeers {
		for _, peer := range device.Peers {
			if peer.PresharedKey.String() == storedPeer.PreSharedKey {
				// Create extended peer and append to list
				ExtendedPeer := types.PeerExtended{
					Peer:          storedPeer,
					TransmitBytes: peer.TransmitBytes,
					ReceiveBytes:  peer.ReceiveBytes,
				}
				ExtendedPeer.LastSeenMillis = peer.LastHandshakeTime.UnixMilli()
				ExtendedPeer.LastIPAddress = peer.Endpoint.IP.String()
				ExtendedPeers = append(ExtendedPeers, ExtendedPeer)
			}
		}
	}

	return ExtendedPeers, nil
}
