package main

import (
	"flag"
	"fmt"
	"net/http"
	"nxy/testsocket/agent/service"
	"os"
	"os/signal"
	"syscall"
)

var (
	webSocketAddr = flag.String("websocket.addr", ":8080", "game agent webSocket address")
	gameServer    = flag.String("grpc.gameserver", "localhost:10020", "game server gRPC server address")
)

func main() {
	flag.Parse()

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	agentService := service.NewAgentService(*gameServer)

	go func() {
		m := http.NewServeMux()
		m.HandleFunc("/", agentService.WebSocketServer)
		errc <- http.ListenAndServe(*webSocketAddr, m)
	}()

	fmt.Println("terminated", <-errc)
}
