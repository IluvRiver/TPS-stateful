package network

import (
	"chat_server_golang/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	SocketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: SocketBufferSize, WriteBufferSize: messageBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}

type Room struct {
	Forward chan *message // 수신되는 메시지를 보관하는 값
	// 들어오는 메시지를 다른 클라이언트들에게 전송을 합니다.

	Join  chan *client // Socket이 연결되는 경우에 작동
	Leave chan *client // Socket이 끊어지는 경우에 대해서 작동

	Clients map[*client]bool // 현재 방에 있는 Client 정보를 저장
	service *service.Service
}

type message struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Room    string    `json:"room"`
	When    time.Time `json:"when"`
}

type client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *client),
		Leave:   make(chan *client),
		Clients: make(map[*client]bool),
		service: service,
	}
}

func (c *client) Read() {
	// 클라이언트가 들어오는 메시지를 읽는 함수
	defer c.Socket.Close()
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			log.Println("READ : ", msg, "client", c.Name)
			log.Println()
			msg.When = time.Now()
			msg.Name = c.Name

			c.Room.Forward <- msg
		}
	}
}

func (c *client) Write() {
	// 클라이언트가 메시지를 전송하는 함수
	defer c.Socket.Close()

	for msg := range c.Send {
		log.Println("WRITE : ", msg, "client", c.Name)
		log.Println()
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true // client가 새로 들어 올떄
		case client := <-r.Leave:
			delete(r.Clients, client) // 나갈 떄에는 map값에서 client를 제거
			close(client.Send)        // 이후 client의 socker을 닫는다.
		case msg := <-r.Forward: // 만약 특정 메시지가 방에 들어오면

			go r.service.InsertChatting(msg.Name, msg.Message, msg.Room)
			//go가 있다면 서브스레드 생성해 백그라운드에서 알아서 진행시키는것

			for client := range r.Clients {
				client.Send <- msg // 모든 client에게 전달 해 준다.
			}
		}
	}
}
func (r *Room) RunInit() {
	// Room에 있는 모든 체널값을을 받는 역할
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) ServeHTTP(c *gin.Context) {
	// 이후 요청이 이렇게 들어오게 된다면 Upgrade를 통해서 소켓을 가져 온다.

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	Socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("---- serveHTTP:", err)
		return
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		log.Fatal("auth cookie is failed", err)
		return
	}

	// 문제가 없다면 client를 생성하여 방에 입장했다고 채널에 전송한다.
	client := &client{
		Socket: Socket,
		Send:   make(chan *message, messageBufferSize),
		Room:   r,
		Name:   authCookie.Value,
	}

	r.Join <- client

	// 또한 defer를 통해서 client가 끝날 떄를 대비하여 퇴장하는 작업을 연기한다.
	defer func() { r.Leave <- client }()

	// 이 후 고루틴을 통해서 write를 실행 시킨다.
	go client.Write()
	// 이 후 메인 루틴에서 read를 실행함으로써 해당 요청을 닫는것을 차단한다.
	// -> 연결을 활성화 시키는 것이다. 채널을 활용하여
	client.Read()
}
