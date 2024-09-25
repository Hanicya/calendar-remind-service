package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

var (
	clients = make(map[*websocket.Conn]bool)
	mu      sync.Mutex
)

// 添加连接
func AddClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	clients[conn] = true
}

// 移除连接
func RemoveClient(conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, conn)
}

// 发送消息
func SendMessage(message string) {
	mu.Lock()
	defer mu.Unlock()
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}
