package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wg-controller/wg-controller/types"
)

type LP_Message struct {
	Topic      string            `json:"topic"`
	Data       string            `json:"data"`
	Attributes map[string]string `json:"attributes"`
	Config     types.Peer        `json:"config,omitempty"`
}

type LP_Client struct {
	Ch           chan LP_Message
	LastConsumed time.Time
}

const ExpiryTime = 60 * time.Second
const PollTimeout = 10 * time.Second

var LP_Clients sync.Map // map[uuid]LP_Client

func InitLongPoll() {
	go func() {
		for {
			GarbageCollectClients()
			time.Sleep(10 * time.Second)
		}
	}()
}

func GarbageCollectClients() {
	LP_Clients.Range(func(key, value interface{}) bool {
		client := value.(*LP_Client)
		if time.Since(client.LastConsumed) > ExpiryTime {
			close(client.Ch)
			LP_Clients.Delete(key)
			log.Println("LongPoll client expired:", key)
		}
		return true
	})
}

func GET_LongPoll(c *gin.Context) {
	// Get the client UUID
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	// Does the client have a channel?
	lpClientInterface, ok := LP_Clients.Load(uuid)
	if !ok {
		// Create a new buffered channel
		ch := make(chan LP_Message, 50)
		lpClient := &LP_Client{
			Ch:           ch,
			LastConsumed: time.Now(),
		}
		LP_Clients.Store(uuid, lpClient)
		lpClientInterface = lpClient
	}

	// Cast the interface once to the pointer
	lpClient := lpClientInterface.(*LP_Client)

	// Send available message or wait
	select {
	case msg := <-lpClient.Ch:
		lpClient.LastConsumed = time.Now()
		c.JSON(200, msg)
		return
	case <-time.After(PollTimeout):
		lpClient.LastConsumed = time.Now()
		c.Status(204) // Tells client to start a new poll
		return
	case <-c.Request.Context().Done():
		return
	}
}

// Sends a message to the long poll client with the given UUID
func SendClientMessage(uuid string, msg LP_Message) error {
	lpClient, ok := LP_Clients.Load(uuid)
	if !ok {
		return errors.New("client not found")
	}
	lpClient.(*LP_Client).Ch <- msg
	return nil
}

func PushPeerConfig(Peer types.Peer) {
	msg := LP_Message{
		Topic:  "peerConfig",
		Config: Peer,
	}
	SendClientMessage(Peer.UUID, msg)
}
