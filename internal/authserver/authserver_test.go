package authserver

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/MaximkaSha/gophkeeper/internal/mockdb"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestAuthGophkeeperServer_UserRegister(t *testing.T) {
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
			data := models.User{
				Email:    "test@test.com",
				Password: "11111",
				Secret:   []byte("secret"),
			}
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().AddUser(gomock.Eq(data))
			Server := AuthGophkeeperServer{
				DB: store,
			}

			s := grpc.NewServer()
			pb.RegisterAuthGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9996")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)

			conn, err := grpc.Dial("localhost:9996", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			c := pb.NewAuthGophkeeperClient(conn)
			_, err = c.UserRegister(context.Background(), &pb.UserRegisterRequest{
				User: data.ToProto(),
			})

			require.NoError(t, err)
			store.EXPECT().AddUser(gomock.Eq(data)).Return(errors.New("no data"))
			_, err = c.UserRegister(context.Background(), &pb.UserRegisterRequest{
				User: data.ToProto(),
			})
			require.Error(t, err)
		})
	}
}

func TestAuthGophkeeperServer_UserLogin(t *testing.T) {
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
			data := models.User{
				Email:    "test@test.com",
				Password: "11111",
				Secret:   []byte("secret"),
			}
			dataHash := models.User{
				Email:    "test@test.com",
				Password: "11111",
				Secret:   []byte("secret"),
			}
			dataHash.HashPassword()
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().GetUser(gomock.Eq(data)).Return(dataHash, nil)
			Server := AuthGophkeeperServer{
				DB: store,
			}

			s := grpc.NewServer()
			pb.RegisterAuthGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9988")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)

			conn, err := grpc.Dial("localhost:9988", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			c := pb.NewAuthGophkeeperClient(conn)
			_, err = c.UserLogin(context.Background(), &pb.UserLoginRequest{
				User: data.ToProto(),
			})

			require.NoError(t, err)
			store.EXPECT().GetUser(gomock.Eq(data)).Return(data, errors.New("no data"))
			_, err = c.UserLogin(context.Background(), &pb.UserLoginRequest{
				User: data.ToProto(),
			})
			require.Error(t, err)
			//wrong pwd
			store.EXPECT().GetUser(gomock.Eq(data)).Return(data, nil)
			_, err = c.UserLogin(context.Background(), &pb.UserLoginRequest{
				User: data.ToProto(),
			})
			require.Error(t, err)
		})
	}
}

func TestAuthGophkeeperServer_JWTClain(t *testing.T) {
	type args struct {
		creds models.User
	}
	tests := []struct {
		name    string
		a       AuthGophkeeperServer
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			a: AuthGophkeeperServer{
				jwtKey: []byte("jwt_key"),
			},
			args: args{
				models.User{
					Email:    "test@test.com",
					Password: "11111111",
					Secret:   []byte("secret"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := tt.a.JWTClain(tt.args.creds)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthGophkeeperServer.JWTClain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			duration := time.Second * 2
			time.Sleep(duration)

			got3, got4, err := tt.a.JWTClain(tt.args.creds)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthGophkeeperServer.JWTClain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotEqual(t, got1, got3)
			require.NotEqual(t, got2, got4)

		})
	}
}

func TestAuthGophkeeperServer_parseToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		a       AuthGophkeeperServer
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			a: AuthGophkeeperServer{
				jwtKey: []byte("secret"),
			},
			args: args{
				token: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := models.User{
				Email:    "t",
				Password: "1",
				Secret:   []byte("s"),
			}
			tt.args.token, _, _ = tt.a.JWTClain(user)
			if err := tt.a.parseToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("AuthGophkeeperServer.parseToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.a.jwtKey = []byte("new_secret")
			err := tt.a.parseToken(tt.args.token)
			require.Error(t, err)
		})
	}
}

func TestAuthGophkeeperServer_Refresh(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.RefreshRequest
	}
	tests := []struct {
		name    string
		a       AuthGophkeeperServer
		args    args
		want    *pb.RefreshResponse
		wantErr bool
	}{
		{
			name: "Test 1",
			a: AuthGophkeeperServer{
				jwtKey: []byte("secret"),
			},
			args: args{
				ctx: context.Background(),
				in: &pb.RefreshRequest{
					Token: &pb.Token{
						Email:   "test@test",
						Token:   "",
						Expires: time.Now().Unix(),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := models.User{
				Email:    "t",
				Password: "1",
				Secret:   []byte("s"),
			}
			tt.args.in.Token.Token, tt.args.in.Token.Expires, _ = tt.a.JWTClain(user)
			_, err := tt.a.Refresh(tt.args.ctx, tt.args.in)
			require.Error(t, err)
			time.Sleep(time.Second * 35)
			resp, err := tt.a.Refresh(tt.args.ctx, tt.args.in)
			require.NoError(t, err)
			require.Equal(t, resp.Token.Token, tt.args.in.Token.Token)
		})
	}
}
