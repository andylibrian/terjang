package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

// NotificationService maintains a collection of subscribers and
// provide a function to broadcast messages to them.
type NotificationService struct {
	subscribers     map[*websocket.Conn]struct{}
	subscribersLock sync.RWMutex
}

// NewNotificationService creates a new notification service.
func NewNotificationService() *NotificationService {
	return &NotificationService{
		subscribers: make(map[*websocket.Conn]struct{}),
	}
}

// AddSubscriber registers a subscriber to the collection.
func (n *NotificationService) AddSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	n.subscribers[conn] = struct{}{}
}

// RemoveSubscriber removes a subscriber from the collection.
func (n *NotificationService) RemoveSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	delete(n.subscribers, conn)
}

// BroadcastMessageToSubscribers sends a message to all of the registered subscribers.
func (n *NotificationService) BroadcastMessageToSubscribers(message []byte) {
	n.subscribersLock.RLock()
	defer n.subscribersLock.RUnlock()

	for conn := range n.subscribers {
		// TODO: conn should be synced
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
