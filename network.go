package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"

	"github.com/vishvananda/netlink"
	"github.com/wg-controller/wg-controller/db"
)

func GetUniqueAddress(usedAddresses []string, serverNetwork string) (string, error) {
	// Parse the network
	_, network, err := net.ParseCIDR(serverNetwork)
	if err != nil {
		return "", err
	}

	// Convert the network IP to a string
	address := network.IP.String()

	// Increment the IP to the first available address
	address, err = nextIP(address, serverNetwork)
	if err != nil {
		return "", err
	}

	// Loop through the addresses in the network
	safetyCounter := 0
	for {
		// Check if the address is used
		used := false
		for _, usedAddress := range usedAddresses {
			if address == usedAddress {
				used = true
				break
			}
		}

		// Return the address if it is not used
		if !used {
			return address, nil
		}

		// Increment the IP
		address, err = nextIP(address, serverNetwork)
		if err != nil {
			return "", err
		}

		// Ensure we don't loop forever
		safetyCounter++
		if safetyCounter > 1000 {
			return "", fmt.Errorf("unable to find a unique address")
		}
	}
}

func nextIP(ipStr string, cidrStr string) (string, error) {
	// Parse the CIDR
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return "", fmt.Errorf("invalid CIDR: %v", err)
	}

	// Parse the IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP: %v", err)
	}

	// Convert the IP to a 4-byte array
	ipBytes := ip.To4()
	if ipBytes == nil {
		return "", fmt.Errorf("invalid IPv4 address")
	}

	// Increment the IP address by 1 (next IP)
	for i := 3; i >= 0; i-- {
		ipBytes[i]++
		if ipBytes[i] != 0 {
			break
		}
	}

	// Ensure the incremented IP is still within the CIDR network
	if ipNet.Contains(ipBytes) {
		return ipBytes.String(), nil
	}
	return "", fmt.Errorf("next IP is outside the CIDR range")
}

func HighestIP(networkCIDR string) (ip string, mask string, err error) {
	// Parse the CIDR
	_, ipNet, err := net.ParseCIDR(networkCIDR)
	if err != nil {
		return "", "", fmt.Errorf("invalid CIDR: %v", err)
	}

	// Get the last IP in the subnet by setting all host bits to 1
	broadcast := ipNet.IP.Mask(ipNet.Mask)
	for i := len(broadcast) - 1; i >= 0; i-- {
		broadcast[i] |= ^ipNet.Mask[i]
	}

	// Subtract 1 to get the highest usable IP (broadcast minus 1)
	for i := len(broadcast) - 1; i >= 0; i-- {
		if broadcast[i] > 0 {
			broadcast[i]--
			break
		}
		broadcast[i] = 255
	}

	// Get the CIDR mask
	maskSize, _ := ipNet.Mask.Size()
	mask = "/" + fmt.Sprintf("%d", maskSize)

	// Return the highest IP address
	return broadcast.String(), mask, nil
}

func GetMask(networkCIDR string) (string, error) {
	// Parse the CIDR
	_, ipNet, err := net.ParseCIDR(networkCIDR)
	if err != nil {
		return "", fmt.Errorf("invalid CIDR: %v", err)
	}

	// Get the CIDR mask
	maskSize, _ := ipNet.Mask.Size()
	mask := "/" + fmt.Sprintf("%d", maskSize)

	// Return the mask
	return mask, nil
}

func SetWireguardInterface() {
	// Set the interface IP address
	cmd1 := exec.Command("ip", "address", "add", ENV.SERVER_ADDRESS, "dev", ENV.WG_INTERFACE)
	err := cmd1.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Set the interface up
	cmd2 := exec.Command("ip", "link", "set", "dev", ENV.WG_INTERFACE, "up")
	err = cmd2.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func InitNetworking() {
	// Create masquerade rule for everything leaving the server
	cmd1 := exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", ENV.EGRESS_INTERFACE, "-j", "MASQUERADE")
	err := cmd1.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func SyncRoutingTable() error {
	// Get all peers from the database
	peers, err := db.GetPeers()
	if err != nil {
		return err
	}

	// Cleanup old routes
	err = CleanupRoutes()
	if err != nil {
		return err
	}

	// Add new routes
	for _, peer := range peers {
		for _, network := range peer.RemoteSubnets {
			err = AddRoute(network, peer.RemoteTunAddress)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CleanupRoutes() error {
	cleanCount := 0
	switch runtime.GOOS {
	case "linux", "darwin":
		routes, _ := netlink.RouteList(nil, 2)
		for _, route := range routes {
			if route.Protocol == 171 {
				err := netlink.RouteDel(&route)
				if err == nil {
					cleanCount++
				}
			}
		}
		log.Println("Cleaned up", cleanCount, "routes")
		return nil
	default:
		return errors.New("unsupported OS")
	}
}

func AddRoute(destination string, gateway string) error {
	switch runtime.GOOS {
	case "linux", "darwin":
		_, dst, err := net.ParseCIDR(destination)
		if err != nil {
			return err
		}
		gw := net.ParseIP(gateway)
		if gw == nil {
			return errors.New("invalid gateway IP")
		}
		route := netlink.Route{
			Dst:      dst,
			Gw:       gw,
			Protocol: 171, // Identifies the route as a WireGuard route
		}
		return netlink.RouteAdd(&route)
	default:
		return errors.New("unsupported OS")
	}
}
