package server

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/MaximkaSha/gophkeeper/internal/mockdb"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGophkeeperServer_AddCipheredData(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			data := models.CipheredData{
				Type: "CC",
				Data: []byte("1"),
				User: "test@test.com",
				Id:   "111-111-111",
			}
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().AddCipheredData(gomock.Eq(data))
			Server := GophkeeperServer{
				DB: store,
			}
			s := grpc.NewServer()
			pb.RegisterGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9999")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)

			conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			c := pb.NewGophkeeperClient(conn)
			_, err = c.AddCipheredData(context.Background(), &pb.AddCipheredDataRequest{
				Data: data.ToProto(),
			})
			require.NoError(t, err)
			store.EXPECT().AddCipheredData(gomock.Eq(data)).Return(errors.New("no data"))
			_, err = c.AddCipheredData(context.Background(), &pb.AddCipheredDataRequest{
				Data: data.ToProto(),
			})
			require.Error(t, err)
		})
	}
}

func TestGophkeeperServer_GetCipheredDataForUserRequest(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			data := models.CipheredData{
				Type: "CC",
				Data: []byte("1"),
				User: "test@test.com",
				Id:   "111-111-111",
			}
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().GetCipheredData(gomock.Eq(data.User)).Return([]models.CipheredData{
				{
					Type: "CC",
					Data: []byte("1"),
					User: "test@test.com",
					Id:   "111-111-111",
				},
				{
					Type: "CC",
					Data: []byte("2"),
					User: "test@test.com",
					Id:   "222-322-222",
				},
			}, nil)
			Server := GophkeeperServer{
				DB: store,
			}
			s := grpc.NewServer()
			pb.RegisterGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9998")
			if err != nil {
				log.Fatal(err)
			}
			_ = listen
			go s.Serve(listen)

			conn, err := grpc.Dial("localhost:9998", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			c := pb.NewGophkeeperClient(conn)
			_, err = c.GetCipheredDataForUserRequest(context.Background(), &pb.GetCipheredDataRequest{
				Email: data.User,
			})
			require.NoError(t, err)
			store.EXPECT().GetCipheredData(gomock.Eq(data.User)).Return([]models.CipheredData{}, errors.New("no data"))
			_, err = c.GetCipheredDataForUserRequest(context.Background(), &pb.GetCipheredDataRequest{
				Email: data.User,
			})
			require.Error(t, err)

		})
	}
}

func TestGophkeeperServer_DelCipheredData(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			data := models.CipheredData{
				Type: "CC",
				Data: []byte("1"),
				User: "test@test.com",
				Id:   "111-111-111",
			}
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().DelCiphereData(gomock.Eq(data.Id)).Return(nil)
			Server := GophkeeperServer{
				DB: store,
			}
			s := grpc.NewServer()
			pb.RegisterGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9997")
			if err != nil {
				log.Fatal(err)
			}
			_ = listen
			go s.Serve(listen)

			conn, err := grpc.Dial("localhost:9997", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			c := pb.NewGophkeeperClient(conn)
			_, err = c.DelCipheredData(context.Background(), &pb.DelCipheredDataRequest{
				Uuid: data.Id,
			})
			require.NoError(t, err)
			store.EXPECT().DelCiphereData(gomock.Eq(data.Id)).Return(errors.New("no data"))
			_, err = c.DelCipheredData(context.Background(), &pb.DelCipheredDataRequest{
				Uuid: data.Id,
			})
			require.Error(t, err)

		})
	}
}
