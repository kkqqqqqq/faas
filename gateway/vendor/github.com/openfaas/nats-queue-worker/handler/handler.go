package handler


import (
	"fmt"
	"log"
	"sync"
)

// CreateNATSQueue ready for asynchronous processing
// type NATSConfig interface {
// 	GetClientID() string
// 	GetMaxReconnect() int
// 	GetReconnectDelay() time.Duration
// }

func CreateNATSQueue(address string, port int, clusterName, channel string, clientConfig NATSConfig) (*NATSQueue, error) {
	var err error
	natsURL := fmt.Sprintf("nats://%s:%d", address, port)
	log.Printf("Opening connection to %s\n", natsURL)

	clientID := clientConfig.GetClientID()

	log.Printf(" natsqueue clientID : %s\n",clientID)

	// If 'channel' is empty, use the previous default.
	if channel == "" {
		channel = "faas-request"
	}

	queue1 := NATSQueue{
		ClientID:       clientID,
		ClusterID:      clusterName,
		NATSURL:        natsURL,
		Topic:          channel,
		maxReconnect:   clientConfig.GetMaxReconnect(),
		reconnectDelay: clientConfig.GetReconnectDelay(),
		ncMutex:        &sync.RWMutex{},
	}

	err = queue1.connect()
	log.Printf("kq:handler.go:natsqueue connected\n" )
	return &queue1, err
}
