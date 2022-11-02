package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MaximkaSha/gophkeeper/internal/authserver"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/MaximkaSha/gophkeeper/internal/server"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

func main() {
	Server := server.NewGophkeeperServer()
	Auth := authserver.NewAuthGophkeeperServer()
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(Auth.AuthFunc)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(Auth.AuthFunc)))
	pb.RegisterGophkeeperServer(s, Server)
	pb.RegisterAuthGophkeeperServer(s, Auth)

	fmt.Println("Gophkeeper Server Started")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
