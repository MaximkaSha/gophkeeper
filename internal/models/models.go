package models

import (
	//"github.com/MaximkaSha/gophkeeper/internal/client"
	"encoding/json"
	"log"

	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CipheredData struct {
	Type string
	Data []byte
	User string
	Id   string
}

func (u *CipheredData) FromProto(proto *pb.CipheredData) {
	u.Data = proto.Data
	u.Id = proto.Uuid
	u.Type = proto.Type.String()
	u.User = proto.Useremail
}

func (u *CipheredData) ToProto() *pb.CipheredData {
	return &pb.CipheredData{
		Data:      u.Data,
		Type:      pb.CipheredData_Type(pb.CipheredData_Type_value[u.Type]),
		Useremail: u.User,
		Uuid:      u.Id,
	}
}

func NewCipheredData(data []byte, email string, dataType string, uuidStr string) *pb.CipheredData {
	return &pb.CipheredData{
		Data:      data,
		Type:      pb.CipheredData_Type(pb.CipheredData_Type_value[dataType]),
		Useremail: email,
		Uuid:      uuidStr,
	}
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Secret   []byte `json:"secret"`
}

func (u *User) FromProto(proto *pb.User) {
	u.Email = proto.Email
	u.Password = proto.Password
	u.Secret = proto.Secret
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Email:    u.Email,
		Password: u.Password,
		Secret:   u.Secret,
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
	return err == nil
}

type Token struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type Password struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Tag      string `json:"tag"`
	ID       string `json:"id"`
}

func (d Password) GetData() []byte {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	data, err := json.Marshal(d)
	if err != nil {
		log.Panic(err)
	}
	return data
}
func (d Password) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}

func (d *Password) SetID() {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
}
func (d Password) Type() string {
	return "PASSWORD"
}

type Data struct {
	Data []byte `json:"data"`
	Tag  string `json:"tag"`
	ID   string `json:"id"`
}

func (d Data) GetData() []byte {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	data, err := json.Marshal(d)
	if err != nil {
		log.Panic(err)
	}
	return data
}
func (d Data) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}
func (d Data) Type() string {
	return "DATA"
}

type Text struct {
	Data string `json:"data"`
	Tag  string `json:"tag"`
	ID   string `json:"id"`
}

func (d Text) GetData() []byte {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	data, err := json.Marshal(d)
	if err != nil {
		log.Panic(err)
	}
	return data
}
func (d Text) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}
func (d Text) Type() string {
	return "TEXT"
}

type CreditCard struct {
	CardNum string `json:"cardnum"`
	Exp     string `json:"exp"`
	Name    string `json:"name"`
	CVV     string `json:"cvv"`
	Tag     string `json:"tag"`
	ID      string `json:"id"`
}

func (d CreditCard) GetData() []byte {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	data, err := json.Marshal(d)
	if err != nil {
		log.Panic(err)
	}
	return data
}
func (d CreditCard) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}
func (d CreditCard) Type() string {
	return "CC"
}

type Storager interface {
	AddUser(User) error
	GetUser(User) (User, error)
	AddCipheredData(CipheredData) error
	GetCipheredData(string) ([]CipheredData, error)
	DelCiphereData(string) error
}

type Dater interface {
	GetData() []byte
	GetID() string
	Type() string
}
