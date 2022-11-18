// Package models implements all data models used by server and client.
package models

import (
	//"github.com/MaximkaSha/gophkeeper/internal/client"
	"encoding/json"
	"log"

	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CipheredData describes ciphered data of all given data types.
type CipheredData struct {
	// Type of data (PASSWORD, CC, TEXT, DATA).
	Type string
	// Marshled and ciphered data.
	Data []byte
	// Email of user which owns data.
	User string
	// Uuid of data.
	ID string
}

// FromProto Function covert data from protobuf to CipheredData.
func (u *CipheredData) FromProto(proto *pb.CipheredData) {
	u.Data = proto.Data
	u.ID = proto.Uuid
	u.Type = proto.Type.String()
	u.User = proto.Useremail
}

// ToProto Function convert CipheredData to protobuff.
func (u *CipheredData) ToProto() *pb.CipheredData {
	return &pb.CipheredData{
		Data:      u.Data,
		Type:      pb.CipheredData_Type(pb.CipheredData_Type_value[u.Type]),
		Useremail: u.User,
		Uuid:      u.ID,
	}
}

// NewCipheredData Construct CipheredData by given values.
func NewCipheredData(data []byte, email string, dataType string, uuidStr string) *pb.CipheredData {
	return &pb.CipheredData{
		Data:      data,
		Type:      pb.CipheredData_Type(pb.CipheredData_Type_value[dataType]),
		Useremail: email,
		Uuid:      uuidStr,
	}
}

// User struct.
type User struct {
	//Email of user.
	Email string `json:"email"`
	// User's password.
	Password string `json:"password"`
	// User's privite key.
	Secret []byte `json:"secret"`
}

// FromProto Covert protobuf User model to models.User.
func (u *User) FromProto(proto *pb.User) {
	u.Email = proto.Email
	u.Password = proto.Password
	u.Secret = proto.Secret
}

// ToProto Convert models.User to protobuf.
func (u *User) ToProto() *pb.User {
	return &pb.User{
		Email:    u.Email,
		Password: u.Password,
		Secret:   u.Secret,
	}
}

// HashPassword Hash users password (bcrypt).
func (u *User) HashPassword() error {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err == nil {
		u.Password = string(passBytes)
		return nil
	}
	return err
}

// CheckPasswordHash Checks users hash.
func (u *User) CheckPasswordHash(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
	return err == nil
}

// Token struct keeps user's jwt token.
type Token struct {
	// Users's email.
	Email string `json:"email"`
	// Users's jwt token.
	Token string `json:"token"`
}

// Password model.
type Password struct {
	// Login.
	Login string `json:"login"`
	// Password.
	Password string `json:"password"`
	// Tag - string which describes for what this password.
	Tag string `json:"tag"`
	// Uniq uuid.
	ID string `json:"id"`
}

// GetData - returns marshled  struct.
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

// GetID - returns uuid of data.
func (d Password) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}

// SetID new uuid for data.
func (d *Password) SetID() {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
}

// Type Return type of data.
func (d Password) Type() string {
	return "PASSWORD"
}

// Data model for files.
type Data struct {
	// Data - slice of file bytes.
	Data []byte `json:"data"`
	// Tag - user given data.
	Tag string `json:"tag"`
	// Uniq uuid.
	ID string `json:"id"`
}

// GetData Return marshled object.
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

// GetID Return uuid of data.
func (d Data) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}

// Type Return type of data.
func (d Data) Type() string {
	return "DATA"
}

// Text Model for texts.
type Text struct {
	// Data - string.
	Data string `json:"data"`
	// User given string.
	Tag string `json:"tag"`
	// Uniq uuid.
	ID string `json:"id"`
}

// GetData Returns marshled object.
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

// GetID Returns uuid of data.
func (d Text) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}

// Type Return type of data.
func (d Text) Type() string {
	return "TEXT"
}

// CreditCard Model of credit cards data.
type CreditCard struct {
	// CardNum - string which contains CC Number.
	CardNum string `json:"cardnum"`
	// Exp - strinf which contains CC expitation date.
	Exp string `json:"exp"`
	// Name - name on credit card.
	Name string `json:"name"`
	// CVV -string contains CVV.
	CVV string `json:"cvv"`
	// User given data.
	Tag string `json:"tag"`
	// Uniq uuid.
	ID string `json:"id"`
}

// GetData Returns marshled object.
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

// GetID Returns uuid of data.
func (d CreditCard) GetID() string {
	if d.ID == "" {
		return uuid.NewString()
	}
	return d.ID
}

// Type Return type of data.
func (d CreditCard) Type() string {
	return "CC"
}

// Storager Interface for database.
type Storager interface {
	AddUser(User) error
	GetUser(User) (User, error)
	AddCipheredData(CipheredData) error
	GetCipheredData(string) ([]CipheredData, error)
	DelCiphereData(string) error
}

// Dater Interface for data converting.
type Dater interface {
	GetData() []byte
	GetID() string
	Type() string
}
