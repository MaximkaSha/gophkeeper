package server

import (
	"context"
	"fmt"
	"log"

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

func (g GophkeeperServer) AddPassword(ctx context.Context, in *pb.AddPasswordRequest) (*pb.AddPasswordResponse, error) {
	var response pb.AddPasswordResponse
	data := models.Password{}
	data.FromProto(in.Password)
	err := g.DB.AddPassword(data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error adding password data`)
	}
	return &response, nil
}

func (g GophkeeperServer) GetPassword(ctx context.Context, in *pb.GetPasswordRequest) (*pb.GetPasswordResponse, error) {
	var response pb.GetPasswordResponse
	data, err := g.DB.GetPassword(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error getting password with uuid: `+in.Id)
	}
	response.Password = data.ToProto()
	return &response, nil
}

func (g GophkeeperServer) DelPassword(ctx context.Context, in *pb.DelPasswordRequest) (*pb.DelPasswordResponse, error) {
	var response pb.DelPasswordResponse
	err := g.DB.DelPassword(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error deliting password with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) UpdatePassword(ctx context.Context, in *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	var response pb.UpdatePasswordResponse
	log.Println(in.Id)
	log.Println(in.Password)
	data := models.Password{}
	data.FromProto(in.Password)
	err := g.DB.UpdatePassword(in.Id, data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error updating password with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) GetAllPassword(ctx context.Context, in *pb.GetAllPasswordRequest) (*pb.GetAllPasswordResponse, error) {
	var response pb.GetAllPasswordResponse

	data, err := g.DB.GetAllPassword()
	if err != nil {
		return &response, status.Errorf(codes.NotFound, `Error getting all passwords`)
	}
	for _, val := range data {
		pVal := val.ToProto()
		response.Password = append(response.Password, pVal)
	}
	return &response, nil
}

func (g GophkeeperServer) AddData(ctx context.Context, in *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	var response pb.AddDataResponse
	log.Println(in.Data)
	return &response, nil
}
func (g GophkeeperServer) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	var response pb.GetDataResponse
	log.Println(in.Id)
	pass := pb.Data{
		Data: []byte("data resp"),
		Tag:  "tag resp",
		Id:   "111",
	}
	response.Data = &pass
	return &response, nil
}

func (g GophkeeperServer) DelData(ctx context.Context, in *pb.DelDataRequest) (*pb.DelDataResponse, error) {
	var response pb.DelDataResponse
	log.Println(in.Id)
	return &response, nil
}

func (g GophkeeperServer) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	var response pb.UpdateDataResponse
	log.Println(in.Data)
	log.Println(in.Id)
	return &response, nil
}

func (g GophkeeperServer) GetAllData(ctx context.Context, in *pb.GetAllDataRequest) (*pb.GetAllDataResponse, error) {
	var response pb.GetAllDataResponse
	for i := 0; i < 10; i++ {
		data := &pb.Data{
			Data: []byte(fmt.Sprint(i)),
			Tag:  fmt.Sprint(i),
			Id:   fmt.Sprint(i),
		}
		response.Data = append(response.Data, data)
	}
	return &response, nil
}

func (g GophkeeperServer) AddText(ctx context.Context, in *pb.AddTextRequest) (*pb.AddTextResponse, error) {
	var response pb.AddTextResponse
	log.Println(in.Text)
	return &response, nil
}

func (g GophkeeperServer) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	var response pb.GetTextResponse
	log.Println(in.Id)
	pass := pb.Text{
		Text: string("text resp"),
		Tag:  "tag resp",
		Id:   "111",
	}
	response.Text = &pass
	return &response, nil
}

func (g GophkeeperServer) DelText(ctx context.Context, in *pb.DelTextRequest) (*pb.DelTextResponse, error) {
	var response pb.DelTextResponse
	log.Println(in.Id)
	return &response, nil
}

func (g GophkeeperServer) UpdateText(ctx context.Context, in *pb.UpdateTextRequest) (*pb.UpdateTextResponse, error) {
	var response pb.UpdateTextResponse
	log.Println(in.Id)
	log.Println(in.Text)
	return &response, nil
}

func (g GophkeeperServer) GetAllText(ctx context.Context, in *pb.GetAllTextRequest) (*pb.GetAllTextResponse, error) {
	var response pb.GetAllTextResponse
	for i := 0; i < 10; i++ {
		data := &pb.Text{
			Text: fmt.Sprint(i),
			Tag:  fmt.Sprint(i),
			Id:   fmt.Sprint(i),
		}
		response.Text = append(response.Text, data)
	}
	return &response, nil
}

func (g GophkeeperServer) AddCreditCard(ctx context.Context, in *pb.AddCreditCardRequest) (*pb.AddCreditCardResponse, error) {
	var response pb.AddCreditCardResponse
	log.Println(in.Creditcard)
	return &response, nil
}

func (g GophkeeperServer) GetCreditCard(ctx context.Context, in *pb.GetCreditCardRequest) (*pb.GetCreditCardResponse, error) {
	var response pb.GetCreditCardResponse
	log.Println(in.Id)
	pass := pb.CreditCard{
		Cardnum: "cc num",
		Exp:     "exp",
		Name:    "name",
		Cvv:     "cvv",
		Tag:     "tag resp",
		Id:      "111",
	}
	response.Creditcard = &pass
	return &response, nil
}

func (g GophkeeperServer) DelCreditCard(ctx context.Context, in *pb.DelCreditCardRequest) (*pb.DelCreditCardResponse, error) {
	var response pb.DelCreditCardResponse
	log.Println(in.Id)
	return &response, nil
}

func (g GophkeeperServer) UpdateCreditCard(ctx context.Context, in *pb.UpdateCreditCardRequest) (*pb.UpdateCreditCardResponse, error) {
	var response pb.UpdateCreditCardResponse
	log.Println(in.Creditcard)
	log.Println(in.Id)
	return &response, nil
}

func (g GophkeeperServer) GetAllCreditCard(ctx context.Context, in *pb.GetAllCreditCardRequest) (*pb.GetAllCreditCardResponse, error) {
	var response pb.GetAllCreditCardResponse
	for i := 0; i < 10; i++ {
		data := &pb.CreditCard{
			Cardnum: fmt.Sprint(i),
			Exp:     fmt.Sprint(i),
			Name:    fmt.Sprint(i),
			Cvv:     fmt.Sprint(i),
			Tag:     fmt.Sprint(i),
			Id:      fmt.Sprint(i),
		}
		response.Creditcard = append(response.Creditcard, data)
	}
	return &response, nil
}
