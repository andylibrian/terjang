package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type NotificationService struct {
	subscribers     map[*websocket.Conn]struct{}
	subscribersLock sync.RWMutex
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		subscribers: make(map[*websocket.Conn]struct{}),
	}
}

func (n *NotificationService) AddSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	n.subscribers[conn] = struct{}{}
}

func (n *NotificationService) RemoveSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	delete(n.subscribers, conn)
}

func (n *NotificationService) BroadcastMessageToSubscribers(message []byte) {
	n.subscribersLock.RLock()
	defer n.subscribersLock.RUnlock()

	for conn := range n.subscribers {
		// TODO: conn should be synced
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
