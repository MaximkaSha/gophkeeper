package server

import (
	"context"

	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/MaximkaSha/gophkeeper/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GophkeeperServer struct {
	DB *storage.Storage
	pb.UnimplementedGophkeeperServer
}

func NewGophkeeperServer() GophkeeperServer {
	return GophkeeperServer{
		DB: storage.NewStorage(),
	}
}

func (g GophkeeperServer) AddCipheredData(ctx context.Context, in *pb.AddCipheredDataRequest) (*pb.AddCipheredDataResponse, error) {
	var response pb.AddCipheredDataResponse
	data := models.CipheredData{}
	data.FromProto(in.Data)
	err := g.DB.AddCipheredData(data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error adding ciphered data `)
	}
	return &response, nil
}

func (g GophkeeperServer) GetCipheredDataForUserRequest(ctx context.Context, in *pb.GetCipheredDataRequest) (*pb.GetCipheredDataResponse, error) {
	var response pb.GetCipheredDataResponse
	user := models.CipheredData{}
	user.FromProto(in.Data)
	data, err := g.DB.GetCipheredData(user)
	if err != nil {
		return &response, status.Errorf(codes.NotFound, `Error getting all ciphered data`)
	}
	for _, val := range data {
		pVal := val.ToProto()
		response.Data = append(response.Data, pVal)
	}
	return &response, nil
}

func (g GophkeeperServer) DelCipheredData(ctx context.Context, in *pb.DelCipheredDataRequest) (*pb.DelCiphereDataResponse, error) {
	var response pb.DelCiphereDataResponse
	data := models.CipheredData{}
	data.FromProto(in.Data)
	err := g.DB.DelCiphereData(data)
	if err != nil {
		return &response, status.Errorf(codes.Unknown, `Error getting all ciphered data`)
	}
	return &response, nil
}
