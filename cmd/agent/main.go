package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	pb "github.com/MaximkaSha/gophkeeper/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func (a *Auth) UnaryAuthClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add the current bearer token to the metadata and call the RPC
	// command
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+a.Token)
	return invoker(ctx, method, req, reply, cc, opts...)
}

type Auth struct {
	Token string
}

func (a *Auth) SetToken(token string) {
	a.Token = token
}

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

	u := pb.NewAuthGophkeeperClient(conn)
	auth.SetToken(TestUser(u))
	c := pb.NewGophkeeperClient(conn)
	TestPassword(c)
	TestData(c)
	TestText(c)
	TestCreditCard(c)
}
func TestUser(c pb.AuthGophkeeperClient) string {
	ctx := context.Background()
	log.Println("------USER--------")

	rand.Seed(time.Now().UnixNano())
	user := pb.User{
		Email:    "test@test.com" + fmt.Sprintf("%v", (rand.Intn(10000))),
		Password: "123456",
	}
	_, err := c.UserRegister(ctx, &pb.UserRegisterRequest{User: &user})
	if err != nil {
		log.Fatal(err)
	}
	response, err := c.UserLogin(ctx, &pb.UserLoginRequest{User: &user})
	if err != nil {
		log.Fatal(err)
	}
	t := time.Unix(response.Token.Expires, 0)
	log.Printf("User:%s, token:%s, expiresAt:%s ", response.Token.Email, response.Token.Token, t.String())
	_, err = c.Refresh(ctx, &pb.RefreshRequest{Token: response.Token})
	if err == nil {
		log.Fatal("no erro. Need token error")
	}
	log.Println("Waiting 35 sec token to expire")
	duration := time.Second * 35
	time.Sleep(duration)
	newToken, err := c.Refresh(ctx, &pb.RefreshRequest{Token: response.Token})
	if err != nil {
		log.Fatalf("Token not refreshed:%v ", err)
	}
	log.Println(newToken)
	return response.Token.Token
}

func Test(c pb.GophkeeperClient) {
	ctx := context.Background()
	log.Println("--------Password-----------")
	user := pb.Password{
		Login:    "Client logn",
		Password: "Client password",
		Tag:      "CLient tag",
	}
	c.AddPassword(ctx, &pb.AddPasswordRequest{
		Password: &user,
	})
	resp, err := c.GetPassword(ctx, &pb.GetPasswordRequest{Id: "Get ID"})
	if err != nil {
		log.Println("err")
	}
	log.Println(resp.Password)
	c.DelPassword(ctx, &pb.DelPasswordRequest{Id: "Del Id"})
	c.UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		Id:       "Update id",
		Password: &user,
	})
	passwords, err := c.GetAllPassword(ctx, &pb.GetAllPasswordRequest{})
	if err != nil {
		log.Println("err")
	}
	log.Println(len(passwords.Password))
	for _, data := range passwords.Password {
		log.Println(data)
	}
	log.Println("--------Data-----------")
	data := pb.Data{
		Data: []byte("client data"),
		Tag:  "CLient tag",
	}
	c.AddData(ctx, &pb.AddDataRequest{
		Data: &data,
	})
	resp1, err := c.GetData(ctx, &pb.GetDataRequest{Id: "Get data ID"})
	if err != nil {
		log.Println("err")
	}
	log.Println(resp1.Data)
	c.DelData(ctx, &pb.DelDataRequest{Id: "Del data Id"})
	c.UpdateData(ctx, &pb.UpdateDataRequest{
		Id:   "Update data id",
		Data: &data,
	})
	datas, err := c.GetAllData(ctx, &pb.GetAllDataRequest{})
	if err != nil {
		log.Println("err")
	}
	for _, data := range datas.Data {
		log.Println(data)
	}
	log.Println("--------Text-----------")
	text := pb.Text{
		Text: "client text",
		Tag:  "CLient tag",
	}
	c.AddText(ctx, &pb.AddTextRequest{
		Text: &text,
	})
	resp2, err := c.GetText(ctx, &pb.GetTextRequest{Id: "Get text ID"})
	if err != nil {
		log.Println("err")
	}
	log.Println(resp2.Text)
	c.DelText(ctx, &pb.DelTextRequest{Id: "Del text Id"})
	c.UpdateText(ctx, &pb.UpdateTextRequest{
		Id:   "Update text id",
		Text: &text,
	})
	texts, err := c.GetAllText(ctx, &pb.GetAllTextRequest{})
	if err != nil {
		log.Println("err")
	}
	for _, data := range texts.Text {
		log.Println(data)
	}
	log.Println("--------CC-----------")
	cc := pb.CreditCard{
		Cardnum: "cc num",
		Exp:     "exp num",
		Name:    "cc name",
		Cvv:     "cvv num",
		Tag:     "CLient tag",
	}
	c.AddCreditCard(ctx, &pb.AddCreditCardRequest{
		Creditcard: &cc,
	})
	resp3, err := c.GetCreditCard(ctx, &pb.GetCreditCardRequest{Id: "Get creditcard ID"})
	if err != nil {
		log.Println("err")
	}
	log.Println(resp3.Creditcard)
	c.DelCreditCard(ctx, &pb.DelCreditCardRequest{Id: "Del cc Id"})
	c.UpdateData(ctx, &pb.UpdateDataRequest{
		Id:   "Update cc id",
		Data: &data,
	})
	ccs, err := c.GetAllCreditCard(ctx, &pb.GetAllCreditCardRequest{})
	if err != nil {
		log.Println("err")
	}
	for _, data := range ccs.Creditcard {
		log.Println(data)
	}
}

func TestPassword(c pb.GophkeeperClient) {
	ctx := context.Background()
	log.Println("--------Password-----------")
	passwords := []*pb.Password{
		{Login: "Max", Password: "Max", Tag: "yandex"},
		{Login: "Max2", Password: "Max2", Tag: "gmail"},
		{Login: "Max3", Password: "max3", Tag: "SZI"},
	}
	for _, pass := range passwords {
		_, err := c.AddPassword(ctx, &pb.AddPasswordRequest{Password: pass})
		if err != nil {
			log.Fatal(err)
		}
	}
	passwordsFromDB, err := c.GetAllPassword(ctx, &pb.GetAllPasswordRequest{})
	if err != nil {
		log.Fatal(err)
	}
	var uuids []string
	log.Println("--------uuids from DB--------")
	for _, pass := range passwordsFromDB.Password {
		log.Println(pass)
		uuids = append(uuids, pass.Id)
	}
	log.Println("--------passwords from DB--------")
	for i, uuid := range uuids {
		pass, err := c.GetPassword(ctx, &pb.GetPasswordRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Password #%v: %s", i, pass)
	}
	updPass := &pb.Password{Login: "NotMax", Password: "NotMax", Tag: "rambler", Id: uuids[0]}
	c.UpdatePassword(ctx, &pb.UpdatePasswordRequest{Password: updPass})
	pass, err := c.GetPassword(ctx, &pb.GetPasswordRequest{Id: uuids[0]})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Password after update: %s", pass)
	for _, uuid := range uuids {
		_, err := c.DelPassword(ctx, &pb.DelPasswordRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("All passwords deleted")
}

func TestData(c pb.GophkeeperClient) {
	ctx := context.Background()
	log.Println("--------Data-----------")
	datas := []*pb.Data{
		{Data: []byte("Max1"), Tag: "yandex data"},
		{Data: []byte("Max2"), Tag: "practikum data"},
		{Data: []byte("Max3"), Tag: "szi data"},
	}
	for _, data := range datas {
		_, err := c.AddData(ctx, &pb.AddDataRequest{Data: data})
		if err != nil {
			log.Fatal(err)
		}
	}
	datasFromDB, err := c.GetAllData(ctx, &pb.GetAllDataRequest{})
	if err != nil {
		log.Fatal(err)
	}
	var uuids []string
	log.Println("--------uuids from DB--------")
	for _, data := range datasFromDB.Data {
		log.Println(data)
		uuids = append(uuids, data.Id)
	}
	log.Println("--------data from DB--------")
	for i, uuid := range uuids {
		data, err := c.GetData(ctx, &pb.GetDataRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Data #%v: %s", i, data)
	}
	updData := &pb.Data{Data: []byte("not data"), Tag: "rambler", Id: uuids[0]}
	c.UpdateData(ctx, &pb.UpdateDataRequest{Data: updData})
	data, err := c.GetData(ctx, &pb.GetDataRequest{Id: uuids[0]})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data after update: %s", data)
	for _, uuid := range uuids {
		_, err := c.DelData(ctx, &pb.DelDataRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("All data deleted")
}

func TestText(c pb.GophkeeperClient) {
	ctx := context.Background()
	log.Println("--------Text-----------")
	texts := []*pb.Text{
		{Text: "Max1", Tag: "yandex text"},
		{Text: "Max2", Tag: "practikum text"},
		{Text: "Max3", Tag: "szi text"},
	}
	for _, text := range texts {
		_, err := c.AddText(ctx, &pb.AddTextRequest{Text: text})
		if err != nil {
			log.Fatal(err)
		}
	}
	textFromDB, err := c.GetAllText(ctx, &pb.GetAllTextRequest{})
	if err != nil {
		log.Fatal(err)
	}
	var uuids []string
	log.Println("--------uuids from DB--------")
	for _, text := range textFromDB.Text {
		log.Println(text)
		uuids = append(uuids, text.Id)
	}
	log.Println("--------text from DB--------")
	for i, uuid := range uuids {
		text, err := c.GetText(ctx, &pb.GetTextRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Text #%v: %s", i, text)
	}
	updText := &pb.Text{Text: "not text", Tag: "rambler", Id: uuids[0]}
	c.UpdateText(ctx, &pb.UpdateTextRequest{Text: updText})
	text, err := c.GetText(ctx, &pb.GetTextRequest{Id: uuids[0]})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Text after update: %s", text)
	for _, uuid := range uuids {
		_, err := c.DelText(ctx, &pb.DelTextRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("All text deleted")
}

func TestCreditCard(c pb.GophkeeperClient) {
	ctx := context.Background()
	log.Println("--------CC-----------")
	ccs := []*pb.CreditCard{
		{Cardnum: "11111", Exp: "11/11", Name: "Max1 Max1", Cvv: "111", Tag: "cc 1"},
		{Cardnum: "22222", Exp: "22/22", Name: "Max2 Max2", Cvv: "222", Tag: "cc 2"},
		{Cardnum: "33333", Exp: "33/33", Name: "Max3 Max3", Cvv: "333", Tag: "cc 3"},
	}
	for _, cc := range ccs {
		_, err := c.AddCreditCard(ctx, &pb.AddCreditCardRequest{Creditcard: cc})
		if err != nil {
			log.Fatal(err)
		}
	}
	ccFromDB, err := c.GetAllCreditCard(ctx, &pb.GetAllCreditCardRequest{})
	if err != nil {
		log.Fatal(err)
	}
	var uuids []string
	log.Println("--------uuids from DB--------")
	for _, cc := range ccFromDB.Creditcard {
		log.Println(cc)
		uuids = append(uuids, cc.Id)
	}
	log.Println("--------cc from DB--------")
	for i, uuid := range uuids {
		cc, err := c.GetCreditCard(ctx, &pb.GetCreditCardRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("CC #%v: %s", i, cc)
	}
	updCC := &pb.CreditCard{Cardnum: "44444", Exp: "44/44", Name: "Max4 Max4", Cvv: "444", Tag: "cc 4", Id: uuids[0]}
	c.UpdateCreditCard(ctx, &pb.UpdateCreditCardRequest{Creditcard: updCC})
	cc, err := c.GetCreditCard(ctx, &pb.GetCreditCardRequest{Id: uuids[0]})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CC after update: %s", cc)
	for _, uuid := range uuids {
		_, err := c.DelCreditCard(ctx, &pb.DelCreditCardRequest{Id: uuid})
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("All cc deleted")
}
