package main

import (
	"context"
	"log"

	pb "github.com/MaximkaSha/gophkeeper/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewGophkeeperClient(conn)
	TestPassword(c)
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
