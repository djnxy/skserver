package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ServerStream interface {
	Send(*Frame) error
	Recv() (*Frame, error)
	CloseSend()
}

type serverStream struct {
	msgChan chan *Frame
	addr    string
}

func (x *serverStream) Send(f *Frame) error {
	resp, err := http.Post(x.addr, "application/json; charset=utf-8", f)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	x.msgChan <- ToFrame(body)
	return nil
}

func (x *serverStream) Recv() (*Frame, error) {
	return <-x.msgChan, nil
}

func (x *serverStream) CloseSend() {
	close(x.msgChan)
}

func NewServerStream(gameServerAddr string) (ServerStream, error) {
	u := url.URL{Scheme: "http", Host: gameServerAddr, Path: "/game"}
	return &serverStream{
		msgChan: make(chan *Frame, 100),
		addr:    u.String(),
	}, nil
}
