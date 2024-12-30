package main

import (
	"errors"
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

func GetWireguardPeer(storedPeer types.Peer) (types.PeerExtended, error) {
	// Get wireguard data
	device, err := wg.Device(ENV.INTERFACE_NAME)
	if err != nil {
		return types.PeerExtended{}, err
	}

	// Find stored peer in wireguard peers
	for _, wgPeer := range device.Peers {
		if wgPeer.PresharedKey.String() == storedPeer.PreSharedKey {
			// Create extended peer
			ExtendedPeer := types.PeerExtended{
				Peer:          storedPeer,
				TransmitBytes: wgPeer.TransmitBytes,
				ReceiveBytes:  wgPeer.ReceiveBytes,
			}
			ExtendedPeer.LastSeenUnixMillis = wgPeer.LastHandshakeTime.UnixMilli()
			ExtendedPeer.LastIPAddress = wgPeer.Endpoint.IP.String()
			return ExtendedPeer, nil
		}
	}

	return types.PeerExtended{}, errors.New("peer not found")
}
