package main

import (
	"log"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/wg-controller/wg-controller/db"
	"github.com/wg-controller/wg-controller/types"
)

var storedPeers sync.Map // Map of peer UUIDs to online status
const minimumPingInterval = 15 * time.Second

func InitAlerts() {
	storedPeers = sync.Map{}
	for {
		// Get all peers
		peers, err := db.GetPeers()
		if err != nil {
			log.Println(err)
			continue
		}

		// Ping each peer in a goroutine
		var wg sync.WaitGroup
		startTime := time.Now()
		for _, peer := range peers {
			wg.Add(1)
			go func(peer types.Peer) {
				defer wg.Done()

				// Check if the peer is online
				online := pingPeer(peer.RemoteTunAddress)

				// Get peer previous status
				p, _ := storedPeers.Load(peer.UUID)
				if p == nil {
					storedPeers.Store(peer.UUID, false)
					return
				}

				// Check if the peer status has changed
				if p.(bool) && !online {
					// Peer has gone offline
					peerStatusAlert(peer.Hostname, false)
				}
				if !p.(bool) && online {
					// Peer has come online
					peerStatusAlert(peer.Hostname, true)
				}

				// Store the new status
				storedPeers.Store(peer.UUID, online)
			}(peer)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Enforce minimum ping interval
		if time.Since(startTime) < minimumPingInterval {
			time.Sleep(minimumPingInterval - time.Since(startTime))
		}
	}
}

func pingPeer(ipAddr string) bool {
	pinger, err := probing.NewPinger(ipAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	pinger.Timeout = 2 * time.Second
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Println(err)
		return false
	}
	stats := pinger.Statistics()

	return stats.PacketsRecv > 0
}

func peerStatusAlert(peerName string, online bool) {
	var event string
	if online {
		event = "🟢 Client Up"
	} else {
		event = "🚨 Client Down"
	}

	var message string
	if online {
		message = peerName + " has come online"
	} else {
		message = peerName + " has gone offline"
	}

	// Send the alert to Slack
	if ENV.SLACK_WEBHOOK != "" {
		msg := NewSlackMessageBody(event, message, "https://"+ENV.PUBLIC_HOST)
		err := SendSlackMessage(ENV.SLACK_WEBHOOK, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func peerCreatedAlert(peerName string) {
	event := "🟢 Client Created"
	message := peerName + " has been created"

	// Send the alert to Slack
	if ENV.SLACK_WEBHOOK != "" {
		msg := NewSlackMessageBody(event, message, "https://"+ENV.PUBLIC_HOST)
		err := SendSlackMessage(ENV.SLACK_WEBHOOK, msg)
		if err != nil {
			log.Println(err)
		}
	}
}
