package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/wg-controller/wg-controller/db"
)

func InitDNS() {
	// Write the upstream DNS server to the dnsmasq configuration
	err := AppendNameserver(ENV.UPSTREAM_DNS)
	if err != nil {
		log.Fatal(err)
	}

	// Sync the DNS configuration
	err = SyncPeersDNS(false)
	if err != nil {
		log.Fatal(err)
	}

	// Get server address without mask
	serverAddress := strings.Split(ENV.SERVER_ADDRESS, "/")[0]

	// Start the DNS server
	startDNS(serverAddress)
}

func startDNS(serverAddress string) {
	log.Println("Starting DNS server...")
	cmd := exec.Command("dnsmasq", "--listen-address="+serverAddress, "--hostsdir=/etc/dnsmasq.d", "--conf-file=/etc/wg-dnsmasq.conf")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Synchronises peers from the database with the dnsmasq configuration
func SyncPeersDNS(restart bool) error {
	log.Println("Syncing hostnames with DNS server")
	// Get all peers from the database
	peers, err := db.GetPeers()
	if err != nil {
		return err
	}

	// Clear the dnsmasq hosts file
	err = ClearDNS()
	if err != nil {
		return err
	}

	// Write the server's IP address to the dnsmasq hosts file
	serverAddress := strings.Split(ENV.SERVER_ADDRESS, "/")[0]
	err = AppendHostname(ENV.SERVER_HOSTNAME, serverAddress)
	if err != nil {
		return err
	}

	// Write the peers to the dnsmasq hosts file
	for _, peer := range peers {
		err = AppendHostname(peer.Hostname, peer.RemoteTunAddress)
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

func RestartDNS() {
	// Get the PID of dnsmasq
	pidCmd := exec.Command("pidof", "dnsmasq")
	pidOutput, err := pidCmd.Output()
	if err != nil {
		log.Fatalf("Failed to get dnsmasq PID: %v", err)
	}

	// Trim any whitespace or newline from the output
	pid := strings.TrimSpace(string(pidOutput))

	// Send SIGHUP to dnsmasq
	killCmd := exec.Command("kill", "-HUP", pid)
	if err := killCmd.Run(); err != nil {
		log.Fatalf("Failed to send SIGHUP to dnsmasq: %v", err)
	}

	log.Println("Successfully sent SIGHUP to dnsmasq")
}

func ClearDNS() error {
	// Check if the file exists
	if _, err := os.Stat("/etc/dnsmasq.d/wg-hosts"); os.IsNotExist(err) {
		// File does not exist, nothing to delete
		return nil
	} else {
		// File exists, delete it
		err := os.Remove("/etc/dnsmasq.d/wg-hosts")
		if err != nil {
			return err
		}
	}

	return nil
}

func AppendNameserver(nameserver string) error {
	// Open/Create the dnsmasq configuration file
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

	// Append the new line
	newEntry := "server=/" + nameserver + "/"
	lines = append(lines, newEntry)

	// Write the lines back to the file
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return err
}

func AppendHostname(hostname string, ip string) error {
	// Open/Create the dnsmasq hosts file
	file, err := os.OpenFile("/etc/dnsmasq.d/wg-hosts", os.O_RDWR|os.O_CREATE, 0644)
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

	// Append the new line
	newEntry := ip + " " + hostname
	lines = append(lines, newEntry)

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
