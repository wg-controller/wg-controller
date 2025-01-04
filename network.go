package main

import (
	"fmt"
	"net"
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
