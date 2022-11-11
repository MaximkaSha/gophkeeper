package main

import (
	"context"
	/*"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/MaximkaSha/gophkeeper/internal/crypto"
	"github.com/MaximkaSha/gophkeeper/internal/models" */
	//pb "github.com/MaximkaSha/gophkeeper/internal/proto"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func (a *Auth) UnaryAuthClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add the current bearer token to the metadata and call the RPC
	// command
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+a.Token)
	return invoker(ctx, method, req, reply, cc, opts...)
}

type Auth struct {
	Token  string
	Email  string
	Hash   string
	Secret []byte
}

func (a *Auth) SetToken(token string, email string, secret []byte) {
	a.Token = token
	a.Email = email
	a.Secret = secret
}

/*
func main() {

	auth := Auth{
		Token: "",
	}
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(auth.UnaryAuthClientInterceptor))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	TestCipher()
	u := pb.NewAuthGophkeeperClient(conn)
	c := pb.NewGophkeeperClient(conn)
	auth.SetToken(TestUser(u))

	TestCipheredData(c, auth)

}

func TestCipher() {
	log.Println("------GOST Cipher Test--------")
	cipher := crypto.NewCrypto([]byte("12345678123456781234567812345678"))
	plain := []byte("some plain text to encrypt with GOST Kuznechik")
	log.Println("Data to encrypt: ", string(plain))
	crypted := cipher.Encrypt(plain)
	log.Println("Decrypted data:")
	log.Println(string(cipher.Decrypt(crypted)))

}

func TestCipheredData(c pb.GophkeeperClient, a Auth) {
	cipher := crypto.NewCrypto(a.Secret)
	ctx := context.Background()
	log.Println("------CipheredData--------")
	testData := []pb.CipheredData{}
	data := models.Password{
		Login:    "Login Test",
		Password: "Password Test",
		Tag:      "Tag Test",
	}
	dataJson, err := json.Marshal(&data)
	dataJson = cipher.Encrypt(dataJson)
	if err != nil {
		log.Fatal(err)
	}
	cData := pb.CipheredData{
		Data:      dataJson,
		Type:      pb.CipheredData_Type(pb.CipheredData_Type_value["PASSWORD"]),
		Useremail: a.Email,
	}
	testData = append(testData, cData)
	data1 := models.Text{
		Data: "text",
		Tag:  "Tag Test",
	}
	data1Json, err := json.Marshal(&data1)
	data1Json = cipher.Encrypt(data1Json)
	if err != nil {
		log.Fatal(err)
	}
	ccData := pb.CipheredData{
		Data:      data1Json,
		Type:      pb.CipheredData_Type(pb.CipheredData_Type_value["TEXT"]),
		Useremail: a.Email,
	}
	testData = append(testData, ccData)
	for _, val := range testData {
		_, err = c.AddCipheredData(ctx, &pb.AddCipheredDataRequest{Data: &val})
		if err != nil {
			log.Fatal(err)
		}
	}
	jData, err := c.GetCipheredDataForUserRequest(ctx, &pb.GetCipheredDataRequest{Email: a.Email})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("-------DATA FROM DB----------")
	for _, val := range jData.Data {
		switch val.Type.String() {
		case "PASSWORD":
			val.Data = cipher.Decrypt(val.Data)
			err = json.Unmarshal(val.Data, &data)
			if err != nil {
				log.Fatal("JSON unmarshal err: ", err)
			}
			log.Println(val.Useremail)
			log.Println(val.Type)
			log.Println(data)
		case "TEXT":
			val.Data = cipher.Decrypt(val.Data)
			err = json.Unmarshal(val.Data, &data1)
			if err != nil {
				log.Fatal("JSON unmarshal err: ", err)
			}
			log.Println(val.Useremail)
			log.Println(val.Type)
			log.Println(data1)
		}

	}
}

func TestUser(c pb.AuthGophkeeperClient) (string, string, []byte) {
	ctx := context.Background()
	log.Println("------USER--------")

	rand.Seed(time.Now().UnixNano())
	user := models.User{
		Email:    "test@test.com" + fmt.Sprintf("%v", (rand.Intn(10000))),
		Password: "123456",
	}
	err := user.HashPassword()
	if err != nil {
		log.Fatal(err)
	}
	userProto := user.ToProto()
	_, err = c.UserRegister(ctx, &pb.UserRegisterRequest{User: userProto})
	if err != nil {
		log.Fatal(err)
	}
	response, err := c.UserLogin(ctx, &pb.UserLoginRequest{User: userProto})
	if err != nil {
		log.Fatal(err)
	}
	t := time.Unix(response.Token.Expires, 0)
	log.Printf("User:%s, token:%s, expiresAt:%s ", response.Token.Email, response.Token.Token, t.String())
	_, err = c.Refresh(ctx, &pb.RefreshRequest{Token: response.Token})
	if err == nil {
		log.Fatal("no erro. Need token error")
	}
	/*
		log.Println("Waiting 35 sec token to expire")
		duration := time.Second * 35
		time.Sleep(duration)
		newToken, err := c.Refresh(ctx, &pb.RefreshRequest{Token: response.Token})
		if err != nil {
			log.Fatalf("Token not refreshed:%v ", err)
		}
		log.Println(newToken)
			return response.Token.Token, response.Token.Email, response.User.Secret
}
*/
