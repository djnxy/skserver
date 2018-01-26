package main

import (
	"flag"
	"fmt"
	"net/http"
	"nxy/testsocket/game/inmem"
	"nxy/testsocket/game/session"
	"os"
	"os/signal"
	"syscall"
)

var gameServer = flag.String("grpc.gameserver", ":10020", "game server gRPC server address")

func main() {
	flag.Parse()
	var users = inmem.NewTestRepository()
	var test session.Service
	test = session.NewService(users)
	mux := http.NewServeMux()
	mux.Handle("/game", session.MakeHandler(test))
	mux.Handle("/game/", session.MakeHandler(test))
	http.Handle("/", mux)
	errc := make(chan error)
	go func() {
		fmt.Println("transport", "http", "address", *gameServer, "msg", "listening")
		errc <- http.ListenAndServe(*gameServer, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("terminated", <-errc)
}
