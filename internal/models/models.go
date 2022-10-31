package models

import (
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
)

type Password struct {
	Login    string
	Password string
	Tag      string
	ID       string
}

func (p *Password) FromProto(proto *pb.Password) {
	p.Login = proto.Login
	p.Password = proto.Password
	p.Tag = proto.Tag
	p.ID = proto.Id
}
func (p *Password) ToProto() *pb.Password {
	return &pb.Password{
		Login:    p.Login,
		Password: p.Password,
		Tag:      p.Tag,
		Id:       p.ID,
	}
}

type Data struct {
	Data []byte
	Tag  string
	ID   string
}

type Text struct {
	Data string
	Tag  string
	ID   string
}

type CreditCard struct {
	CardNum string
	Exp     string
	Name    string
	CVV     string
	ID      string
}

type Storager interface {
	AddPassword(Password) error
	GetPassword(string) (Password, error)
	DelPassword(string) error
	UpdatePassword(string, Password) error
	GetAllPassword() ([]Password, error)

	AddData(Data) error
	GetData(string) (Data, error)
	DelData(string) error
	UpdateData(string, Data) error
	GetAllData() ([]Data, error)

	AddText(Text) error
	GetText(string) (Text, error)
	DelText(string) error
	UpdateText(string, Text) error
	GetAllText() ([]Text, error)

	AddCreditCard(CreditCard) error
	GetCreditCard(string) (CreditCard, error)
	DelCreditCard(string) error
	UpdateCreditCard(string, CreditCard)
	GetAllCreditCard() ([]CreditCard, error)
}
