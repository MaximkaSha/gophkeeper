// Package server is backend for the client.
// It implements gRPC endpoints and DB handling.
// Configured via json file.
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
	"google.golang.org/grpc/credentials"
)

func main() {
	Server := server.NewGophkeeperServer()
	Auth := authserver.NewAuthGophkeeperServer()
	listen, err := net.Listen("tcp", Server.Config.Addr)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile(Server.Config.CertFile, Server.Config.CertKey)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(Auth.AuthFunc)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(Auth.AuthFunc)),
		grpc.Creds(creds))
	pb.RegisterGophkeeperServer(s, Server)
	pb.RegisterAuthGophkeeperServer(s, Auth)

	fmt.Println("Gophkeeper Server Started")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
