package citest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Test_citest(t *testing.T) {
	log.Println("-----Starting gRPC server-----")
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
	go s.Serve(listen)
	log.Println("-----Gophkeeper Server Started-----")

	client := client.NewClient("Test", "Test")
	log.Println("-----Gophkeeper Client Started-----")
	testCI(t, client)
}

func testCI(t *testing.T, client *client.Client) {
	ctx := context.Background()
	log.Println("------USER Register/Login--------")
	rand.Seed(time.Now().UnixNano())
	user := models.User{
		Email:    "test@test.com" + fmt.Sprintf("%v", (rand.Intn(10000))),
		Password: "123456",
	}
	err := client.UserRegister(ctx, user)
	require.NoError(t, err, "User Register Error")
	err = client.UserLogin(ctx, user)
	require.NoError(t, err, "User Login Error")
	log.Println("----Register/Login completed with no erros----")
	go client.RefreshToken(ctx)
	log.Println("------CipheredData Add--------")
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
	log.Println("------CipheredData Added successfully--------")
	err = client.GetAllDataFromDB(ctx)
	require.NoError(t, err, "Error getting data from DB")
	log.Print("-------Checking data from DB----------")
	for _, val := range client.AllData {
		switch val.Type {
		case "PASSWORD":
			pass := models.Password{}
			json.Unmarshal(val.JData, &pass)
			require.Equal(t, dataP.Login, pass.Login, "Logins not equal")
			require.Equal(t, dataP.Password, pass.Password, "Passwords not equal")
			require.Equal(t, dataP.Tag, pass.Tag, "Tag not equal")
			log.Println("Passwords equal!")
		case "CC":
			cc := models.CreditCard{}
			json.Unmarshal(val.JData, &cc)
			require.Equal(t, dataC.CardNum, cc.CardNum, "CC Num not equal")
			require.Equal(t, dataC.Exp, cc.Exp, "CC Exp not equal")
			require.Equal(t, dataC.CVV, cc.CVV, "CC CVV not equal")
			require.Equal(t, dataC.Name, cc.Name, "CC Name not equal")
			require.Equal(t, dataC.Tag, cc.Tag, "CC Tag not equal")
			log.Println("CC equal!")
		case "DATA":
			data := models.Data{}
			json.Unmarshal(val.JData, &data)
			require.Equal(t, dataD.Data, data.Data, "Data not equal")
			require.Equal(t, dataD.Tag, data.Tag, "Data tag not equal")
			log.Println("DATA equal!")
		case "TEXT":
			text := models.Text{}
			json.Unmarshal(val.JData, &text)
			require.Equal(t, dataT.Data, text.Data, "Text data not equal")
			require.Equal(t, dataT.Tag, text.Tag, "Text tag not equal")
			log.Println("TEXT equal!")
		}

	}
	log.Print("-------All data from DB checked----------")
	log.Print("-------Deleting data from DB----------")
	for _, val := range client.AllData {
		err := client.DelData(ctx, val.ID)
		require.NoError(t, err, "Error delelting data from DB")
	}
	log.Print("-------Deleting data from DB----------")
	log.Print("-------All test passed----------")
}
