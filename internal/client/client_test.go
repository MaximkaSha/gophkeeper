package client

import (
	"context"
	"errors"
	"log"
	"net"
	"reflect"
	"testing"

	"github.com/MaximkaSha/gophkeeper/internal/authserver"
	"github.com/MaximkaSha/gophkeeper/internal/crypto"
	"github.com/MaximkaSha/gophkeeper/internal/mockdb"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/MaximkaSha/gophkeeper/internal/server"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient_AddData(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
	}{
		{
			name: "test",
			c:    &Client{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			data := models.Password{
				Login:    "test",
				Password: "test",
				Tag:      "test",
				ID:       "test",
			}
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().AddCipheredData(gomock.Any()).Times(2)
			Server := server.GophkeeperServer{
				DB: store,
			}
			s := grpc.NewServer()
			pb.RegisterGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9994")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)
			conn, err := grpc.Dial("localhost:9994", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			tt.c.serverClient = pb.NewGophkeeperClient(conn)
			tt.c.crypto = *crypto.NewCrypto([]byte("12345678123456781234567812345678"))
			err = tt.c.AddData(context.Background(), data)
			require.NoError(t, err)
			err = tt.c.AddData(context.Background(), data)
			require.NoError(t, err)
			store.EXPECT().AddCipheredData(gomock.Any()).Return(errors.New("no data"))
			err = tt.c.AddData(context.Background(), data)
			require.Error(t, err)

		})
	}
}

func TestClient_DelData(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
	}{
		{
			name: "test",
			c:    &Client{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			data := models.Password{
				Login:    "test",
				Password: "test",
				Tag:      "test",
				ID:       "test",
			}
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().DelCiphereData(gomock.Any()).Times(2)
			Server := server.GophkeeperServer{
				DB: store,
			}
			s := grpc.NewServer()
			pb.RegisterGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9993")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)
			conn, err := grpc.Dial("localhost:9993", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			tt.c.serverClient = pb.NewGophkeeperClient(conn)
			tt.c.crypto = *crypto.NewCrypto([]byte("12345678123456781234567812345678"))
			allData := []AllData{
				{
					ID: "000-111-00",
				},
				{
					ID: "test",
				},
				{
					ID: "111-111-111",
				},
			}
			tt.c.AllData = allData
			err = tt.c.DelData(context.Background(), data.ID)
			require.NoError(t, err)
			allData = []AllData{
				{
					ID: "test",
				},
			}
			tt.c.AllData = allData
			err = tt.c.DelData(context.Background(), data.ID)
			require.NoError(t, err)
			store.EXPECT().DelCiphereData(gomock.Any()).Return(errors.New("no data"))
			err = tt.c.DelData(context.Background(), data.ID)
			require.Error(t, err)

		})
	}
}

func TestClient_GetAllDataFromDB(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
	}{
		{
			name: "test",
			c:    &Client{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().GetCipheredData(gomock.Any())
			Server := server.GophkeeperServer{
				DB: store,
			}
			s := grpc.NewServer()
			pb.RegisterGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9992")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)
			conn, err := grpc.Dial("localhost:9992", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			tt.c.serverClient = pb.NewGophkeeperClient(conn)
			tt.c.crypto = *crypto.NewCrypto([]byte("12345678123456781234567812345678"))
			err = tt.c.GetAllDataFromDB(context.Background())
			require.NoError(t, err)
			store.EXPECT().GetCipheredData(gomock.Any()).Return([]models.CipheredData{}, errors.New("no data"))
			err = tt.c.GetAllDataFromDB(context.Background())
			require.Error(t, err)
			store.EXPECT().GetCipheredData(gomock.Any()).Return([]models.CipheredData{
				{
					Type: "CC",
					Data: []byte("testtesttesttest"),
					User: "test@tst.com",
					ID:   "111-111-1111",
				},
				{
					Type: "CC",
					Data: []byte("testtesttesttest"),
					User: "test@tst.com",
					ID:   "222-111-1111",
				},
				{
					Type: "CC",
					Data: []byte("testtesttesttest"),
					User: "test@tst.com",
					ID:   "444-111-1111",
				},
			}, nil)
			err = tt.c.GetAllDataFromDB(context.Background())
			require.Error(t, err)
		})
	}
}

func TestClient_AddDataToLocalStorageUI(t *testing.T) {
	type args struct {
		ctx context.Context
		v   any
	}
	tests := []struct {
		name string
		c    *Client
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.AddDataToLocalStorageUI(tt.args.ctx, tt.args.v)
		})
	}
}

/*
func TestClient_RefreshToken(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
	}{
		{
			name: "test",
			c: &Client{
				auth: &Auth{
					Email: "test@test.com",
					Token: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := authserver.AuthGophkeeperServer{}
			auth.SetJWTKey("123456578")
			listen, err := net.Listen("tcp", "localhost:9991")
			if err != nil {
				log.Fatal(err)
			}
			s := grpc.NewServer()
			pb.RegisterAuthGophkeeperServer(s, auth)
			go s.Serve(listen)
			conn, err := grpc.Dial("localhost:9991", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			tt.c.authClient = pb.NewAuthGophkeeperClient(conn)
			token := "wrong token"
			tt.c.auth.Token = token
			tt.c.RefreshToken(context.Background())
			require.Equal(t, token, tt.c.auth.Token)

		})
	}
} */

func TestClient_UserRegister(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
	}{
		{
			name: "test",
			c: &Client{
				auth: &Auth{
					Email:  "test@test.com",
					Token:  "",
					Secret: []byte("1234567812345678123456781234567"),
				},
			},
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
				Secret:   []byte("12345678123456781234567812345678"),
			}
			dataHash.HashPassword()
			store := mockdb.NewMockStorager(ctrl)
			Server := authserver.AuthGophkeeperServer{
				DB: store,
			}

			s := grpc.NewServer()
			pb.RegisterAuthGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9989")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)

			conn, err := grpc.Dial("localhost:9989", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			tt.c.authClient = pb.NewAuthGophkeeperClient(conn)
			store.EXPECT().AddUser(gomock.Any()).Return(nil)
			err = tt.c.UserRegister(context.Background(), data)
			require.NoError(t, err)
			store.EXPECT().AddUser(gomock.Any()).Return(errors.New("no data"))
			err = tt.c.UserRegister(context.Background(), data)
			require.Error(t, err)
		})
	}
}

func TestClient_UserLogin(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
	}{
		{
			name: "test",
			c: &Client{
				auth: &Auth{
					Email:  "test@test.com",
					Token:  "",
					Secret: []byte("1234567812345678123456781234567"),
				},
			},
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
				Secret:   []byte("12345678123456781234567812345678"),
			}
			dataHash.HashPassword()
			store := mockdb.NewMockStorager(ctrl)
			store.EXPECT().GetUser(gomock.Eq(data)).Return(dataHash, nil)
			Server := authserver.AuthGophkeeperServer{
				DB: store,
			}

			s := grpc.NewServer()
			pb.RegisterAuthGophkeeperServer(s, Server)
			listen, err := net.Listen("tcp", "localhost:9995")
			if err != nil {
				log.Fatal(err)
			}
			go s.Serve(listen)

			conn, err := grpc.Dial("localhost:9995", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			tt.c.authClient = pb.NewAuthGophkeeperClient(conn)
			tt.c.UserLogin(context.Background(), data)
			require.NoError(t, err)
			store.EXPECT().GetUser(gomock.Eq(data)).Return(data, errors.New("no data"))
			err = tt.c.UserLogin(context.Background(), data)
			require.Error(t, err)
		})
	}
}

func TestClient_UnmarshalProtoData(t *testing.T) {
	type args struct {
		val *pb.CipheredData
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "CC",
			c: &Client{
				crypto: *crypto.NewCrypto([]byte("12345678123456781234567812345678")),
			},
			args: args{
				val: &pb.CipheredData{
					Data:      []byte("test cc test cc test cc test cc="),
					Type:      pb.CipheredData_CC,
					Useremail: "mail",
					Uuid:      "111-111-1111",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DATA",
			c: &Client{
				crypto: *crypto.NewCrypto([]byte("12345678123456781234567812345678")),
			},
			args: args{
				val: &pb.CipheredData{
					Data:      []byte("test cc test cc test cc test cc="),
					Type:      pb.CipheredData_DATA,
					Useremail: "mail",
					Uuid:      "111-111-1111",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "PASSWORD",
			c: &Client{
				crypto: *crypto.NewCrypto([]byte("12345678123456781234567812345678")),
			},
			args: args{
				val: &pb.CipheredData{
					Data:      []byte("test cc test cc test cc test cc="),
					Type:      pb.CipheredData_PASSWORD,
					Useremail: "mail",
					Uuid:      "111-111-1111",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "TEXT",
			c: &Client{
				crypto: *crypto.NewCrypto([]byte("12345678123456781234567812345678")),
			},
			args: args{
				val: &pb.CipheredData{
					Data:      []byte("test cc test cc test cc test cc="),
					Type:      pb.CipheredData_TEXT,
					Useremail: "mail",
					Uuid:      "111-111-1111",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.UnmarshalProtoData(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UnmarshalProtoData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UnmarshalProtoData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_PrinStorage(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
	}{
		{
			name: "test",
			c: &Client{
				LocalStorage: NewLocalStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.PrinStorage()
		})
	}
}

func TestLocalStorage_DelFromLocalStorageUI(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		l    *LocalStorage
		args args
	}{
		{
			name: "data",
			l:    NewLocalStorage(),
			args: args{
				models.Data{
					Data: []byte("ss"),
					Tag:  "11",
					ID:   "data",
				},
			},
		},
		{
			name: "cc",
			l:    NewLocalStorage(),
			args: args{
				models.CreditCard{
					CardNum: "",
					Exp:     "",
					Name:    "",
					CVV:     "",
					Tag:     "11",
					ID:      "data",
				},
			},
		},
		{
			name: "text",
			l:    NewLocalStorage(),
			args: args{
				models.Text{
					Data: "text",
					Tag:  "11",
					ID:   "data",
				},
			},
		},
		{
			name: "password",
			l:    NewLocalStorage(),
			args: args{
				models.Password{
					Password: "111",
					Login:    "qqqq",
					Tag:      "11",
					ID:       "data",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.AppendOrUpdate(tt.args.v)
			tt.l.AppendOrUpdate(tt.args.v)
			tt.l.DelFromLocalStorageUI(tt.args.v)
		})
	}
}
