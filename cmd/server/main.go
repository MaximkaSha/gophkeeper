package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/MaximkaSha/gophkeeper/internal/server"
	"google.golang.org/grpc"
)

func main() {
	Server := server.NewGophkeeperServer()

	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterGophkeeperServer(s, Server)
	fmt.Println("Сервер gRPC начал работу")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
