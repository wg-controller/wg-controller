package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/wg-controller/wg-controller/db"
)

var dnsKillChan chan bool

func InitDNS() {
	// Sync the DNS configuration
	err := SyncPeersDNS(false)
	if err != nil {
		log.Fatal(err)
	}

	// Get server address without mask
	serverAddress := strings.Split(ENV.SERVER_ADDRESS, "/")[0]

	// Start the DNS server in a goroutine
	go startDNS(serverAddress)
}

func startDNS(serverAddress string) {
	log.Println("Starting DNS server...")
	cmd := exec.Command("dnsmasq", "--listen-address="+serverAddress, "--conf-file=/etc/wg-dnsmasq.conf", "-k")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create a channel to kill the DNS server
	dnsKillChan = make(chan bool)

	// Listen for kill signal
	select {
	case <-dnsKillChan:
		// Kill the dnsmasq process if the channel receives a signal
		log.Println("Stopping DNS server...")
		cmd.Process.Kill()
	}
}

func RestartDNS() {
	// Send a signal to the DNS server to restart
	select {
	case dnsKillChan <- true:
	default:
		log.Println("DNS server is not running")
	}

	go startDNS(strings.Split(ENV.SERVER_ADDRESS, "/")[0])
}

// Synchronises peers from the database with the dnsmasq configuration
func SyncPeersDNS(restart bool) error {
	log.Println("Syncing hostnames with DNS server")
	// Get all peers from the database
	peers, err := db.GetPeers()
	if err != nil {
		return err
	}

	// Clear the dnsmasq configuration
	err = ClearDNS()
	if err != nil {
		return err
	}

	// Write the server's IP address to the dnsmasq configuration
	serverAddress := strings.Split(ENV.SERVER_ADDRESS, "/")[0]
	err = AppendHostname(ENV.SERVER_HOSTNAME, serverAddress)
	if err != nil {
		return err
	}

	// Write the upstream DNS server to the dnsmasq configuration
	err = AppendNameserver(ENV.UPSTREAM_DNS)
	if err != nil {
		return err
	}

	// Write the peers to the dnsmasq configuration
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

func ClearDNS() error {
	// Open the dnsmasq configuration file
	file, err := os.OpenFile("/etc/wg-dnsmasq.conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Truncate the file
	err = file.Truncate(0)
	return err
}

func AppendNameserver(nameserver string) error {
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

	// Append the new line
	newEntry := "address=/" + hostname + "/" + ip
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
