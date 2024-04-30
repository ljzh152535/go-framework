package websocketLog

import (
	"golang.org/x/net/websocket"
	"log"
)

// websocket客户端
type client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

// 客户端管理
type clientManager struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

var manager = clientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
}

func (manager *clientManager) start() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[seelog] error:%+v", err)
		}
	}()

	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (c *client) write() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteClose(1)
				return
			}
			c.socket.Write(message)
		}
	}
}

func (c *client) read() {
	/*defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		manager.broadcast <- jsonMessage
	}*/
}
