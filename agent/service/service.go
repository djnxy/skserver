package service

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/satori/go.uuid"

	"github.com/gorilla/websocket"
)

const (
	WebSocketReadDeadline int = 15 //秒
)

var (
	ErrReadWebSocket  = errors.New("err on read data from webSocket")
	ErrWriteWebSocket = errors.New("err on write data to webSocket")
	ErrReadRPCStream  = errors.New("err on read data from rpc stream")
	ErrWriteRPCStream = errors.New("err on write data to rpc stream")
)

type Session struct {
	id      string //会话随机ID
	as      *AgentService
	conn    *websocket.Conn //webSocket连接
	stream  ServerStream
	wg      sync.WaitGroup
	dieChan chan struct{}
}

func NewSession(id string, conn *websocket.Conn, stream ServerStream, service *AgentService) *Session {
	s := &Session{
		id:      id,
		as:      service,
		conn:    conn,
		stream:  stream,
		dieChan: make(chan struct{}),
	}
	return s
}

func (s *Session) ForwardToClient(payload []byte) (err error) {
	err = s.conn.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		fmt.Println("err", err.Error(), "msg", "webSocket write err")
	}
	return
}

func (s *Session) ForwardToServer(c *ClientMessage) (err error) {
	//转发到服务器
	err = s.stream.Send(ClientToFrame(c, s.id))
	fmt.Println("id", s.id, "msg", "forward client data to server")
	if err != nil {
		fmt.Println("err", err.Error(), "msg", "stream send err")
	}
	return
}

func (s *Session) HandleClient() {
	errChan := make(chan error)
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		newclientmsg := &ClientMessage{Cmd: "session", Data: map[string]interface{}{"sid": s.id}, Errno: "0"}
		err := s.ForwardToClient(newclientmsg.ToByte())
		if err != nil {
			errChan <- ErrWriteWebSocket
			return
		}
		//处理来自客户端的数据
		for {
			//s.conn.SetReadDeadline(time.Now().Add(time.Duration(WebSocketReadDeadline) * time.Second))
			_, bytes, err := s.conn.ReadMessage()
			if err != nil {
				fmt.Println("id", s.id, "err", err.Error(), "msg", "session webSocket read data err")
				errChan <- ErrReadWebSocket
				return
			}
			err = s.ForwardToServer(ToClientMessage(bytes))
			if err != nil {
				errChan <- ErrWriteRPCStream
				return
			}
			select {
			case <-s.dieChan:
				return
			default:
			}

		}
	}()

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		newclientmsg := &ClientMessage{Cmd: "login", Data: "", Errno: "0"}
		err := s.ForwardToServer(newclientmsg)
		if err != nil {
			errChan <- ErrWriteRPCStream
			return
		}
		//处理来自服务器的数据
		for {
			f, err := s.stream.Recv()
			if err != nil {
				fmt.Println("err", err.Error(), "msg", "stream recv err")
				errChan <- ErrReadRPCStream
				return
			}
			if f == nil {
				return
			}
			//err = s.ForwardToClient(f.ToByte())
			err = s.as.dispatcher(f)
			if err != nil {
				errChan <- ErrWriteWebSocket
				return
			}
			select {
			case <-s.dieChan:
				return
			default:

			}
		}

	}()
	s.as.closeSession(s.id, <-errChan)
}

type AgentService struct {
	mtx            sync.RWMutex
	SessDieChan    chan int64
	Sessions       map[string]*Session
	gameServerAddr string
}

func NewAgentService(gameaddr string) AgentService {
	fmt.Println("msg", "start new agent service")

	agent := AgentService{
		Sessions:       make(map[string]*Session),
		gameServerAddr: gameaddr,
	}
	return agent
}

func (a AgentService) WebSocketServer(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("err", err.Error(), "msg", "err on create webSocket for session")
		return
	}
	stream, err := NewServerStream(a.gameServerAddr)
	if err != nil {
		conn.Close()
		fmt.Println("err", err.Error(), "msg", "err on create rpc stream for session")
		return
	}
	a.mtx.Lock()
	defer a.mtx.Unlock()
	randid, _ := uuid.NewV4()
	id := randid.String()
	a.Sessions[id] = NewSession(id, conn, stream, &a)
	go a.Sessions[id].HandleClient()
	fmt.Println("start new session")
}

func (a AgentService) dispatcher(f *Frame) error {
	if len(f.Sessionid) <= 0 {
		return errors.New("empty ids in Frame")
	}
	for _, sid := range f.Sessionid {
		ss, exist := a.Sessions[sid]
		if !exist {
			continue
		}
		ss.ForwardToClient(FrameToClient(f).ToByte())
	}
	return nil
}

func (a AgentService) closeSession(id string, err error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	bf := len(a.Sessions)
	ss := a.Sessions[id]
	lastclientmsg := &ClientMessage{Cmd: "logout", Data: "", Errno: "0"}
	ss.ForwardToServer(lastclientmsg)
	ss.conn.Close()
	ss.stream.CloseSend()
	close(ss.dieChan)
	ss.wg.Wait()
	delete(a.Sessions, id)
	af := len(a.Sessions)
	fmt.Println("id", id, "before", bf, "after", af, "err", err.Error(), "msg", "close session")
}
