// Packager server implements gRPC server.
package server

import (
	"context"

	"github.com/MaximkaSha/gophkeeper/internal/config"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/MaximkaSha/gophkeeper/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GophkeeperServer - server struct.
type GophkeeperServer struct {
	// Database interface.
	DB models.Storager
	pb.UnimplementedGophkeeperServer
	// Serves's config.
	Config *config.ServerConfig
}

// COnstructor of GophkeeperServer
func NewGophkeeperServer() GophkeeperServer {
	config := config.NewServerConfig()
	return GophkeeperServer{
		DB:     storage.NewStorage(config.DSN),
		Config: config,
	}
}

// AddCipheredData - gRPC endpoint push data to server's db.
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

// GetCipheredDataForUserRequest gRPC endpoint returns all data for given user.
func (g GophkeeperServer) GetCipheredDataForUserRequest(ctx context.Context, in *pb.GetCipheredDataRequest) (*pb.GetCipheredDataResponse, error) {
	var response pb.GetCipheredDataResponse
	data, err := g.DB.GetCipheredData(in.Email)
	if err != nil {
		return &response, status.Errorf(codes.NotFound, `Error getting all ciphered data`)
	}
	for _, val := range data {
		pVal := val.ToProto()
		response.Data = append(response.Data, pVal)
	}
	return &response, nil
}

// DelCipheredData - gRPC endpoint delete data from by given uuid
func (g GophkeeperServer) DelCipheredData(ctx context.Context, in *pb.DelCipheredDataRequest) (*pb.DelCiphereDataResponse, error) {
	var response pb.DelCiphereDataResponse
	err := g.DB.DelCiphereData(in.Uuid)
	if err != nil {
		return &response, status.Errorf(codes.Unknown, `Error getting all ciphered data`)
	}
	return &response, nil
}
