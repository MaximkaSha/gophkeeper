package models

import (
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string
	Password string
}

func (u *User) FromProto(proto *pb.User) {
	u.Email = proto.Email
	u.Password = proto.Password
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *User) HashPassword() error {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err == nil {
		u.Password = string(passBytes)
		return nil
	}
	return err
}

func (u *User) CheckPasswordHash(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
	return !(err == nil)
}

type Token struct {
	Email string
	Token string
}

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

func (d *Data) FromProto(proto *pb.Data) {
	d.Data = proto.Data
	d.Tag = proto.Tag
	d.ID = proto.Id
}

func (d *Data) ToProto() *pb.Data {
	return &pb.Data{
		Data: d.Data,
		Tag:  d.Tag,
		Id:   d.ID,
	}
}

type Text struct {
	Data string
	Tag  string
	ID   string
}

func (d *Text) FromProto(proto *pb.Text) {
	d.Data = proto.Text
	d.Tag = proto.Tag
	d.ID = proto.Id
}

func (d *Text) ToProto() *pb.Text {
	return &pb.Text{
		Text: d.Data,
		Tag:  d.Tag,
		Id:   d.ID,
	}
}

type CreditCard struct {
	CardNum string
	Exp     string
	Name    string
	CVV     string
	Tag     string
	ID      string
}

func (c *CreditCard) FromProto(proto *pb.CreditCard) {
	c.CardNum = proto.Cardnum
	c.Exp = proto.Exp
	c.Name = proto.Name
	c.CVV = proto.Cvv
	c.ID = proto.Id
	c.Tag = proto.Tag
}

func (c *CreditCard) ToProto() *pb.CreditCard {
	return &pb.CreditCard{
		Cardnum: c.CardNum,
		Exp:     c.Exp,
		Name:    c.Name,
		Cvv:     c.CVV,
		Tag:     c.Tag,
		Id:      c.ID,
	}
}

type Storager interface {
	AddUser(User) error
	GetUser(User) (User, error)

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
