package server

import (
	"context"

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
	data := models.Data{}
	data.FromProto(in.Data)
	err := g.DB.AddData(data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error adding  data`)
	}
	return &response, nil
}
func (g GophkeeperServer) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	var response pb.GetDataResponse
	data, err := g.DB.GetData(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error getting data with uuid: `+in.Id)
	}
	response.Data = data.ToProto()
	return &response, nil
}

func (g GophkeeperServer) DelData(ctx context.Context, in *pb.DelDataRequest) (*pb.DelDataResponse, error) {
	var response pb.DelDataResponse
	err := g.DB.DelData(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error deliting data with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	var response pb.UpdateDataResponse
	data := models.Data{}
	data.FromProto(in.Data)
	err := g.DB.UpdateData(in.Id, data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error updating data with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) GetAllData(ctx context.Context, in *pb.GetAllDataRequest) (*pb.GetAllDataResponse, error) {
	var response pb.GetAllDataResponse
	data, err := g.DB.GetAllData()
	if err != nil {
		return &response, status.Errorf(codes.NotFound, `Error getting all data`)
	}
	for _, val := range data {
		pVal := val.ToProto()
		response.Data = append(response.Data, pVal)
	}
	return &response, nil
}

func (g GophkeeperServer) AddText(ctx context.Context, in *pb.AddTextRequest) (*pb.AddTextResponse, error) {
	var response pb.AddTextResponse
	data := models.Text{}
	data.FromProto(in.Text)
	err := g.DB.AddText(data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error adding text `)
	}
	return &response, nil
}

func (g GophkeeperServer) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	var response pb.GetTextResponse
	data, err := g.DB.GetText(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error getting text with uuid: `+in.Id)
	}
	response.Text = data.ToProto()
	return &response, nil
}

func (g GophkeeperServer) DelText(ctx context.Context, in *pb.DelTextRequest) (*pb.DelTextResponse, error) {
	var response pb.DelTextResponse
	err := g.DB.DelText(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error deliting text with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) UpdateText(ctx context.Context, in *pb.UpdateTextRequest) (*pb.UpdateTextResponse, error) {
	var response pb.UpdateTextResponse
	data := models.Text{}
	data.FromProto(in.Text)
	err := g.DB.UpdateText(in.Id, data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error updating text with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) GetAllText(ctx context.Context, in *pb.GetAllTextRequest) (*pb.GetAllTextResponse, error) {
	var response pb.GetAllTextResponse
	data, err := g.DB.GetAllText()
	if err != nil {
		return &response, status.Errorf(codes.NotFound, `Error getting all text`)
	}
	for _, val := range data {
		pVal := val.ToProto()
		response.Text = append(response.Text, pVal)
	}
	return &response, nil
}

func (g GophkeeperServer) AddCreditCard(ctx context.Context, in *pb.AddCreditCardRequest) (*pb.AddCreditCardResponse, error) {
	var response pb.AddCreditCardResponse
	data := models.CreditCard{}
	data.FromProto(in.Creditcard)
	err := g.DB.AddCreditCard(data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error adding cc `)
	}
	return &response, nil
}

func (g GophkeeperServer) GetCreditCard(ctx context.Context, in *pb.GetCreditCardRequest) (*pb.GetCreditCardResponse, error) {
	var response pb.GetCreditCardResponse
	data, err := g.DB.GetCreditCard(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error getting cc with uuid: `+in.Id)
	}
	response.Creditcard = data.ToProto()
	return &response, nil
}

func (g GophkeeperServer) DelCreditCard(ctx context.Context, in *pb.DelCreditCardRequest) (*pb.DelCreditCardResponse, error) {
	var response pb.DelCreditCardResponse
	err := g.DB.DelCreditCard(in.Id)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error deliting cc with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) UpdateCreditCard(ctx context.Context, in *pb.UpdateCreditCardRequest) (*pb.UpdateCreditCardResponse, error) {
	var response pb.UpdateCreditCardResponse
	data := models.CreditCard{}
	data.FromProto(in.Creditcard)
	err := g.DB.UpdateCreditCard(in.Id, data)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, `Error updating cc with uuid: `+in.Id)
	}
	return &response, nil
}

func (g GophkeeperServer) GetAllCreditCard(ctx context.Context, in *pb.GetAllCreditCardRequest) (*pb.GetAllCreditCardResponse, error) {
	var response pb.GetAllCreditCardResponse
	data, err := g.DB.GetAllCreditCard()
	if err != nil {
		return &response, status.Errorf(codes.NotFound, `Error getting all cc`)
	}
	for _, val := range data {
		pVal := val.ToProto()
		response.Creditcard = append(response.Creditcard, pVal)
	}
	return &response, nil
}

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

func (g GophkeeperServer) GetCipheredDataForUserRequest(ctx context.Context, in *pb.GetCipheredDataRequest) (*pb.GetCipheredDataResponse, error) {
	var response pb.GetCipheredDataResponse
	user := models.CipheredData{}
	user.FromProto(in.Data)
	data, err := g.DB.GetCipheredData(user)
	if err != nil {
		return &response, status.Errorf(codes.NotFound, `Error getting all ciphered data`)
	}
	for _, val := range data {
		pVal := val.ToProto()
		response.Data = append(response.Data, pVal)
	}
	return &response, nil
}
