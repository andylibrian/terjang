package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

/* NotificationService is a struct type that has 2 fields;
* map[*websocket.Conn]struct{} subscribers
* sync.RWMutex subscribersLock
****************************************************************/
type NotificationService struct {
	subscribers     map[*websocket.Conn]struct{}
	subscribersLock sync.RWMutex
}

/* NewNotificationService is a function of NotificationService pointer
*******************************************************************/
func NewNotificationService() *NotificationService {
	return &NotificationService{
		subscribers: make(map[*websocket.Conn]struct{}),
	}
}

/* AddSubscriber is a method to add subscriber that has a receiver type of *NotificationService
* and takes a parameter of type conn
******************************************************************/
func (n *NotificationService) AddSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	n.subscribers[conn] = struct{}{}
}

/* RemoveSubscriber is a method to remove a subscriber that has a receiver type of *NotificationService
* and takes a parameter of type conn
******************************************************************/
func (n *NotificationService) RemoveSubscriber(conn *websocket.Conn) {
	n.subscribersLock.Lock()
	defer n.subscribersLock.Unlock()

	delete(n.subscribers, conn)
}

/* BroadcastMessageToSubscribers is a method to writeMessage to a websocket that has a receiver type of *NotificationService
* and takes a parameter of type []byte
******************************************************************/
func (n *NotificationService) BroadcastMessageToSubscribers(message []byte) {
	n.subscribersLock.RLock()
	defer n.subscribersLock.RUnlock()

	for conn := range n.subscribers {
		// TODO: conn should be synced
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
