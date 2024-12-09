package svc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// send buffer size
	bufSize = 256
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type WsClient struct {
	logx.Logger
	svcCtx *ServiceContext
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan types.WsMessage
	id   string
	// operator for websocket
	wsoperator *WsOperator
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *WsClient) readPump() {
	defer func() {
		c.svcCtx.Hub.unregister <- c
		c.conn.Close()
		c.Logger.Infof("Websocket client %s closed", c.conn.RemoteAddr().String())
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		c.Logger.Infof("Recv websocket message: %s", string(message))
		if err != nil {
			logx.Error(err)
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// }
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.processMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Logger.Error("The hub closed the channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.Logger.Error(err)
				return
			}
			messageBytes, _ := json.Marshal(msg)
			w.Write(messageBytes)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				mBytes, _ := json.Marshal(<-c.send)
				w.Write(mBytes)
			}

			if err := w.Close(); err != nil {
				c.Logger.Error(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Logger.Error(err)
				return
			}
		}
	}
}

func (c *WsClient) processMessage(message []byte) {
	var msg types.WsMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.Logger.Errorf("Failed to parse websocket message: %v", err)
		return
	}
	switch msg.MessageType {
	case constant.WS_MSG_TYPE_PODLOG:
		c.wsoperator.GetPodLogs(msg)
	default:
		c.svcCtx.Hub.broadcast <- msg
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(svcCtx *ServiceContext, w http.ResponseWriter, r *http.Request) {
	logging := logx.WithContext(r.Context())
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error(err)
		return
	}
	id := r.FormValue("id")
	logging.Infof("WebSocket received a connection from: %s, id=%s", conn.RemoteAddr().String(), id)
	client := &WsClient{
		Logger:     logging,
		svcCtx:     svcCtx,
		conn:       conn,
		send:       make(chan types.WsMessage, bufSize),
		id:         id,
		wsoperator: NewWsOperator(r.Context(), svcCtx),
	}
	client.svcCtx.Hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
