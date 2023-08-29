package websocket

import (
	"context"
	"encoding/json"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
	"github.com/winc-link/hummingbird/internal/pkg/i18n"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/middleware"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsClient struct {
	id   string
	hub  *WsServer
	ctx  *gin.Context
	lc   logger.LoggingClient
	dic  *di.Container
	conn *websocket.Conn
	send chan WsResponse
}

type WsServer struct {
	lc logger.LoggingClient

	// Registered clients.
	clients map[*wsClient]bool

	clientIdMap map[string]*wsClient

	// Register requests from the clients.
	register chan *wsClient

	// Unregister requests from clients.
	unregister chan *wsClient

	broadcast chan WsResponse
	ctx       context.Context
	dic       *di.Container
}

func NewServer(dic *di.Container) *WsServer {
	var lc = container.LoggingClientFrom(dic.Get)
	c := &WsServer{
		lc:          lc,
		clients:     make(map[*wsClient]bool),
		clientIdMap: make(map[string]*wsClient),
		register:    make(chan *wsClient),
		unregister:  make(chan *wsClient),
		broadcast:   make(chan WsResponse),
		ctx:         context.Background(),
		dic:         dic,
	}
	go c.run()
	go c.listenBroadcast()

	return c
}

func (s *WsServer) Handle(c *gin.Context) {
	lc := s.lc
	w := c.Writer
	r := c.Request
	s.ctx = r.Context()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		lc.Error("upgrade webSocket err", err)
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}

	// 对于client，以ip:userid 作为唯一请求方，用于断线重连
	clientId := strings.Split(conn.RemoteAddr().String(), ":")[0]
	if r.Header.Get("X-Real-Ip") != "" {
		clientId = r.Header.Get("X-Real-Ip")
	}
	lc.Debugf("ws client ip: %s", clientId)

	value, ok := c.Get(constants.JwtParsedInfo)
	if ok {
		claim, ok := value.(*middleware.CustomClaims)
		if ok {
			clientId = clientId + ":" + strconv.Itoa(int(claim.ID))
		}
	}

	client := &wsClient{
		id:   clientId,
		hub:  s,
		conn: conn,
		ctx:  c,
		dic:  s.dic,
		lc:   lc,
		send: make(chan WsResponse),
	}
	s.register <- client

	go client.writePump()
	go client.readPump()
}

func (s *WsServer) run() {
	for {
		select {
		//TODO: 缺少 done 时的退出
		case client := <-s.register:
			s.clients[client] = true
			s.clientIdMap[client.id] = client
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
			}
			if _, ok := s.clientIdMap[client.id]; ok {
				delete(s.clientIdMap, client.id)
			}
		case data := <-s.broadcast:
			s.lc.Debugf("broadcast forward message to alertclient: %v, data: %+v", len(s.clients), data)
			for client := range s.clients {
				select {
				case client.send <- data:
				default:
				}
			}
		}
	}
}

// 监听内部业务推送前端广播消息
func (s *WsServer) listenBroadcast() {
	sc := container.StreamClientFrom(s.dic.Get)
	for {
		select {
		case data := <-sc.Recv():
			var resp httphelper.CommonResponse
			if data.ErrCode == errort.DefaultSuccess {
				resp = httphelper.NewSuccessCommonResponse(data.Data)
			} else {
				resp = httphelper.NewFailWithI18nResponse(s.ctx, errort.NewCommonErr(data.ErrCode, nil))
			}
			s.broadcast <- WsResponse{
				Code: data.Code,
				Data: resp,
			}
		}
	}
}

func (c *wsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !message.Data.Success {
				if message.Data.ErrorCode != errort.ContainerRunFail { //如果是ContainerRunFail错误，把原错误返回出去方便排查问题。
					message.Data.ErrorMsg = i18n.TransCode(c.ctx, message.Data.ErrorCode, nil)
				}
			}

			messageBody, _ := json.Marshal(message)

			c.lc.Infof("websocket to resp data: %s", string(messageBody))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.lc.Warn("client send channel closed!")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.lc.Errorf("websocket NextWriter err:", err)
				return
			}
			_, err = w.Write(messageBody)
			if err != nil {
				c.lc.Errorf("websocket Write err:", err)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *wsClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				c.lc.Errorf("ReadMessage close info: %v", err)
			}
			break
		}
		if !json.Valid(msg) {
			c.lc.Errorf("ReadMessage data not is json format")
			continue
		}
		c.lc.Infof("websocket req %+v", string(msg))
		d := WsData{}
		err = json.Unmarshal(msg, &d)
		if err != nil {
			c.lc.Errorf("ReadMessage data unmarshal err: %v", err)
			continue
		}
		if f, ok := wsFuncMap[d.Code]; ok {
			go f(c, d.Data, d.Code)
		}
	}
}

func (c *wsClient) sendData(code dtos.WsCode, d httphelper.CommonResponse) {
	resData := WsResponse{}
	resData.Code = code
	resData.Data = d

	c.send <- resData
}

// 切换语言
func (c *wsClient) ChangeLang(lang string) {
	c.ctx.Set(constants.AcceptLanguage, lang)
}
