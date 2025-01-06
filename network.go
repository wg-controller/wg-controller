package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
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
	// Get the highest IP address in the server's CIDR range
	highestIP, mask, err := HighestIP(ENV.SERVER_CIDR)
	if err != nil {
		log.Fatal(err)
	}

	// Set the interface IP address
	cmd1 := exec.Command("ip", "address", "add", highestIP+mask, "dev", ENV.WG_INTERFACE)
	err = cmd1.Run()
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
