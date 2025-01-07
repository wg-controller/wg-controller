package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lampy255/wg-controller/db"
)

func InitDNS() {
	// Sync the DNS configuration
	err := SyncPeersDNS(false)
	if err != nil {
		log.Fatal(err)
	}

	// Get server address without mask
	serverAddress := strings.Split(ENV.SERVER_ADDRESS, "/")[0]

	// Start the DNS server
	log.Println("Starting DNS server...")
	cmd := exec.Command("dnsmasq", "--listen-address="+serverAddress, "--conf-file=/etc/wg-dnsmasq.conf")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err.Error() + ":" + string(out))
	}
	log.Println("DNS server started")
}

func RestartDNS() {
	cmd := exec.Command("kill", "-HUP", "$(pidof dnsmasq)")
	cmd.Run()
}

// Synchronises peers from the database with the dnsmasq configuration
func SyncPeersDNS(restart bool) error {
	log.Println("Syncing hostnames with DNS server")
	// Get all peers from the database
	peers, err := db.GetPeers()
	if err != nil {
		return err
	}

	// Write the server's IP address to the dnsmasq configuration
	serverAddress := strings.Split(ENV.SERVER_ADDRESS, "/")[0]
	err = SyncDNSEntry(ENV.SERVER_HOSTNAME, serverAddress, false)
	if err != nil {
		return err
	}

	// Write the peers to the dnsmasq configuration
	for _, peer := range peers {
		err = SyncDNSEntry(peer.Hostname, peer.RemoteTunAddress, false)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// Restart the DNS server
	if restart {
		RestartDNS()
	}

	return nil
}

func SyncDNSEntry(hostname string, ip string, delete bool) error {
	// Open the dnsmasq configuration file
	file, err := os.OpenFile("/etc/wg-dnsmasq.conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Split the file into lines
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Remove any existing line with the same hostname
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], hostname) {
			// Remove the line
			lines = append(lines[:i], lines[i+1:]...)
			i-- // Decrease the index to re-check the current position
		}
	}

	// Append the new line
	if !delete {
		newEntry := "address=/" + hostname + "/" + ip
		lines = append(lines, newEntry)
	}

	// Write the lines back to the file
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
