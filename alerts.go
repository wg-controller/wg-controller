package main

import (
	"log"
	"time"

	"github.com/wg-controller/wg-controller/db"
	"github.com/wg-controller/wg-controller/types"
)

var storedPeers []types.Peer

func InitAlerts() {
	for {
		time.Sleep(15 * time.Second)

		// Get all peers
		peers, err := getAllPeers()
		if err != nil {
			log.Println(err)
			continue
		}

		// Check if the peers have changed
		for _, peer := range peers {
			for _, storedPeer := range storedPeers {
				if peer.UUID == storedPeer.UUID {
					// Determine peer status by last seen timestamp
					onlineNow := time.Since(time.UnixMilli(peer.LastSeenUnixMillis)).Seconds() < 60
					onlineBefore := time.Since(time.UnixMilli(storedPeer.LastSeenUnixMillis)).Seconds() < 60

					if onlineNow && !onlineBefore {
						// Peer has come online
						peerStatusAlert(peer.Hostname, true)
					}
					if !onlineNow && onlineBefore {
						// Peer has gone offline
						peerStatusAlert(peer.Hostname, false)
					}

					break
				}
			}
		}

		// Update the stored peers
		storedPeers = peers
	}
}

func getAllPeers() ([]types.Peer, error) {
	peers, err := db.GetPeers()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var extendedPeers []types.Peer
	for _, peer := range peers {
		if peer.Enabled {
			extendedPeer, err := GetWireguardPeer(peer)
			if err != nil {
				log.Println(err)
				continue // Skip this peer
			}
			extendedPeers = append(extendedPeers, extendedPeer)
		} else {
			extendedPeers = append(extendedPeers, peer)
		}
	}

	return extendedPeers, nil
}

func peerStatusAlert(peerName string, online bool) {
	var event string
	if online {
		event = "Client Up 🟢"
	} else {
		event = "Client Down 🚨"
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
	event := "Client Created 🟢"
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
