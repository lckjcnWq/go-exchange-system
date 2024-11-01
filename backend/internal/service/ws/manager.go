package ws

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Manager struct {
	clients   *gmap.Map   // 所有客户端连接
	broadcast chan []byte // 广播消息通道
	upgrader  websocket.Upgrader
	sync.RWMutex
}

var manager *Manager

func Init() {
	manager = &Manager{
		clients:   gmap.New(true),
		broadcast: make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true //允许所有来源
			},
		},
	}
	go manager.handleBroadcast()
}

func (m *Manager) handleBroadcast() {
	for message := range m.broadcast {
		m.clients.Iterator(func(k interface{}, v interface{}) bool {
			client := v.(*websocket.Conn)
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				g.Log().Error(nil, "Error sending message to client:", err)
				client.Close()
				m.clients.Remove(k)
			}
			return true
		})
	}
}

func (m *Manager) HandleConnection(r *ghttp.Request) {
	conn, err := m.upgrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		g.Log().Error(nil, "Failed to websocket upgrade connection:", err)
		return
	}
	//保存连接
	clientId := r.Get("userId").String()
	m.clients.Set(clientId, conn)

	defer func() {
		conn.Close()
		m.clients.Remove(clientId)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				g.Log().Error(r.Context(), "WebSocket read error:", err)
			}
			break
		}
		// 处理接收到的消息
		m.handleMessage(msg)
	}
}

func (m *Manager) handleMessage(msg []byte) {
	// 处理接收到的消息
	g.Log().Debug(nil, "Received message:", string(msg))
}

// Broadcast 广播消息给所有客户端
func (m *Manager) Broadcast(message []byte) {
	m.broadcast <- message
}

// Get 获取WebSocket管理器实例
func Get() *Manager {
	return manager
}
