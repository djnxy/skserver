package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

var agentAddr = flag.String("addr", "localhost:8080", "service address")

func main() {
	flag.Parse()
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	u := url.URL{Scheme: "ws", Host: *agentAddr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
	}
	defer c.Close()
	fmt.Println("socket connected")
	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				errc <- err
				return
			}
			fmt.Println("read: ", string(message))
		}
	}()
	go func() {
		var msg string
		for {
			fmt.Scanln(&msg)
			fmt.Println("input:", msg)
			err := c.WriteMessage(websocket.BinaryMessage, []byte(msg))
			if err != nil {
				fmt.Println("input:", err)
				errc <- err
				return
			}
		}
	}()
	fmt.Println("terminated", <-errc)
}
