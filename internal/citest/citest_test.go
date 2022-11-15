package citest

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/MaximkaSha/gophkeeper/internal/authserver"
	"github.com/MaximkaSha/gophkeeper/internal/client"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/MaximkaSha/gophkeeper/internal/server"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Test_citest(t *testing.T) {
	tt := suite.Suite{}
	tt.T().Log("-----Starting gRPC server-----")
	Server := server.NewGophkeeperServer()
	Auth := authserver.NewAuthGophkeeperServer()
	listen, err := net.Listen("tcp", Server.Config.Addr)
	require.NoError(t, err, "error starting gRPC server")

	creds, err := credentials.NewServerTLSFromFile(Server.Config.CertFile, Server.Config.CertKey)
	require.NoError(t, err, "error  reading server ceritficate")
	s := grpc.NewServer(grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(Auth.AuthFunc)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(Auth.AuthFunc)),
		grpc.Creds(creds))
	pb.RegisterGophkeeperServer(s, Server)
	pb.RegisterAuthGophkeeperServer(s, Auth)
	go s.Serve(listen)
	tt.T().Log("-----Gophkeeper Server Started-----")
	client := client.NewClient("Test", "Test")
	tt.T().Log("-----Gophkeeper Client Started-----")
	testCI(t, client)
}

func testCI(t *testing.T, client *client.Client) {
	tt := suite.Suite{}
	ctx := context.Background()
	tt.T().Log("------USER Register/Login--------")
	rand.Seed(time.Now().UnixNano())
	user := models.User{
		Email:    "test@test.com" + fmt.Sprintf("%v", (rand.Intn(10000))),
		Password: "123456",
	}
	err := client.UserRegister(ctx, user)
	require.NoError(t, err, "User Register Error")
	err = client.UserLogin(ctx, user)
	require.NoError(t, err, "User Login Error")
	tt.T().Log("----Register/Login completed with no erros----")
	go client.RefreshToken(ctx)
	tt.T().Log("------CipheredData Add--------")
	testData := []models.Dater{}
	dataP := models.Password{
		Login:    "Login Test",
		Password: "Password Test",
		Tag:      "Tag Test",
	}
	testData = append(testData, dataP)
	dataT := models.Text{
		Data: "text",
		Tag:  "Tag Test",
	}
	testData = append(testData, dataT)
	dataD := models.Data{
		Data: []byte("some secret data"),
		Tag:  "Data test",
	}
	testData = append(testData, dataD)
	dataC := models.CreditCard{
		CardNum: "test cc num",
		Exp:     "test cc exp",
		Name:    "test cc name",
		CVV:     "test cc cvv",
		Tag:     "CC test",
	}
	testData = append(testData, dataC)
	for _, val := range testData {
		err = client.AddData(ctx, val)
		require.NoError(t, err, "Error adding data to storage.")
	}
	tt.T().Log("------CipheredData Added successfully--------")
	err = client.GetAllDataFromDB(ctx)
	require.NoError(t, err, "Error getting data from DB")
	tt.T().Log("-------Checking data from DB----------")
	for _, val := range client.AllData {
		switch val.Type {
		case "PASSWORD":
			pass := models.Password{}
			json.Unmarshal(val.JData, &pass)
			require.Equal(t, dataP.Login, pass.Login, "Logins not equal")
			require.Equal(t, dataP.Password, pass.Password, "Passwords not equal")
			require.Equal(t, dataP.Tag, pass.Tag, "Tag not equal")
			tt.T().Log("		Passwords equal!")
		case "CC":
			cc := models.CreditCard{}
			json.Unmarshal(val.JData, &cc)
			require.Equal(t, dataC.CardNum, cc.CardNum, "CC Num not equal")
			require.Equal(t, dataC.Exp, cc.Exp, "CC Exp not equal")
			require.Equal(t, dataC.CVV, cc.CVV, "CC CVV not equal")
			require.Equal(t, dataC.Name, cc.Name, "CC Name not equal")
			require.Equal(t, dataC.Tag, cc.Tag, "CC Tag not equal")
			tt.T().Log("		CC equal!")
		case "DATA":
			data := models.Data{}
			json.Unmarshal(val.JData, &data)
			require.Equal(t, dataD.Data, data.Data, "Data not equal")
			require.Equal(t, dataD.Tag, data.Tag, "Data tag not equal")
			tt.T().Log("		DATA equal!")
		case "TEXT":
			text := models.Text{}
			json.Unmarshal(val.JData, &text)
			require.Equal(t, dataT.Data, text.Data, "Text data not equal")
			require.Equal(t, dataT.Tag, text.Tag, "Text tag not equal")
			tt.T().Log("		TEXT equal!")
		}

	}
	tt.T().Log("-------All data from DB checked----------")
	tt.T().Log("-------Deleting data from DB----------")
	for _, val := range client.AllData {
		err := client.DelData(ctx, val.ID)
		require.NoError(t, err, "Error delelting data from DB")
	}
	tt.T().Log("-------Deleting data from DB----------")
	tt.T().Log("-------All test passed----------")
}
