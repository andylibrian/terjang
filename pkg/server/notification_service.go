package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

// NotificationService is ...
type NotificationService struct {
	subscribers     map[*websocket.Conn]struct{}
	subscribersLock sync.RWMutex
}

// NewNotificationService is ...
func NewNotificationService() *NotificationService {
	return &NotificationService{
		subscribers: make(map[*websocket.Conn]struct{}),
	}
}

// AddSubscriber is ...
func (n *NotificationService) AddSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	n.subscribers[conn] = struct{}{}
}

// RemoveSubscriber is ...
func (n *NotificationService) RemoveSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	delete(n.subscribers, conn)
}

// BroadcastMessageToSubscribers ..
func (n *NotificationService) BroadcastMessageToSubscribers(message []byte) {
	n.subscribersLock.RLock()
	defer n.subscribersLock.RUnlock()

	for conn := range n.subscribers {
		// TODO: conn should be synced
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
