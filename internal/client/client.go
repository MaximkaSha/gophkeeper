package client

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/MaximkaSha/gophkeeper/internal/config"
	"github.com/MaximkaSha/gophkeeper/internal/crypto"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func (a *Auth) UnaryAuthClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+a.Token)
	return invoker(ctx, method, req, reply, cc, opts...)
}

type Auth struct {
	Token  string
	Email  string
	Hash   string
	Secret []byte
}

type LocalStorage struct {
	DataStorage     []models.Data
	TextStorage     []models.Text
	CCStorage       []models.CreditCard
	PasswordStorage []models.Password
}

func (l *LocalStorage) AppendOrUpdate(v any) {
	switch v := v.(type) {
	case models.Data:
		for i := range l.DataStorage {
			if l.DataStorage[i].ID == v.ID {
				l.DataStorage[i].Data = v.Data
				l.DataStorage[i].Tag = v.Tag
				return
			}
		}
		l.DataStorage = append(l.DataStorage, v)
	case models.Password:
		for i := range l.PasswordStorage {
			if l.PasswordStorage[i].ID == v.ID {
				l.PasswordStorage[i].Login = v.Login
				l.PasswordStorage[i].Password = v.Password
				l.PasswordStorage[i].Tag = v.Tag
				return
			}
		}
		l.PasswordStorage = append(l.PasswordStorage, v)
	case models.CreditCard:
		for i := range l.CCStorage {
			if l.CCStorage[i].ID == v.ID {
				l.CCStorage[i].CardNum = v.CardNum
				l.CCStorage[i].Exp = v.Exp
				l.CCStorage[i].CVV = v.CVV
				l.CCStorage[i].Name = v.Name
				l.CCStorage[i].Tag = v.Tag
				return
			}
		}
		l.CCStorage = append(l.CCStorage, v)
	case models.Text:
		for i := range l.TextStorage {
			if l.TextStorage[i].ID == v.ID {
				l.TextStorage[i].Data = v.Data
				l.TextStorage[i].Tag = v.Tag
				return
			}
		}
		l.TextStorage = append(l.TextStorage, v)
	}

}
func (l *LocalStorage) DelFromLocalStorage(v any) {
	switch v := v.(type) {
	case models.Data:
		for i := range l.DataStorage {
			if l.DataStorage[i].ID == v.ID {
				ret := make([]models.Data, 0)
				l.DataStorage = append(ret, l.DataStorage[:i]...)
				return
			}
		}
	case models.Password:
		for i := range l.PasswordStorage {
			if l.PasswordStorage[i].ID == v.ID {
				ret := make([]models.Password, 0)
				l.PasswordStorage = append(ret, l.PasswordStorage[:i]...)
				return
			}
		}
	case models.CreditCard:
		for i := range l.CCStorage {
			if l.CCStorage[i].ID == v.ID {
				ret := make([]models.CreditCard, 0)
				l.CCStorage = append(ret, l.CCStorage[:i]...)
				return
			}
		}
	case models.Text:
		for i := range l.TextStorage {
			if l.TextStorage[i].ID == v.ID {
				ret := make([]models.Text, 0)
				l.TextStorage = append(ret, l.TextStorage[:i]...)
				return
			}
		}
	}

}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		DataStorage:     []models.Data{},
		TextStorage:     []models.Text{},
		CCStorage:       []models.CreditCard{},
		PasswordStorage: []models.Password{},
	}
}

type Client struct {
	authClient   pb.AuthGophkeeperClient
	serverClient pb.GophkeeperClient
	crypto       crypto.Crypto
	currentUser  models.User
	auth         *Auth
	LocalStorage *LocalStorage
	Config       *config.ClientConfig
	BuildVersion string
	BuildTime    string
}

func NewClient(bv string, bt string) *Client {
	auth := &Auth{
		Token: "",
	}
	config := config.NewClientConfig()
	credsTmp, err := credentials.NewClientTLSFromFile(config.CertFile, "")
	if err != nil {
		log.Fatalf("loading GRPC key error: %s", err.Error())
	}
	conn, err := grpc.Dial(config.Addr, grpc.WithTransportCredentials(credsTmp),
		grpc.WithUnaryInterceptor(auth.UnaryAuthClientInterceptor))
	if err != nil {
		log.Fatal(err)
	}
	u := pb.NewAuthGophkeeperClient(conn)
	c := pb.NewGophkeeperClient(conn)
	return &Client{
		authClient:   u,
		serverClient: c,
		auth:         auth,
		LocalStorage: NewLocalStorage(),
		BuildVersion: bv,
		BuildTime:    bt,
	}
}

func (c *Client) PrinStorage() {
	log.Println("CC Storage: ", c.LocalStorage.CCStorage)
	log.Println("Password Storage: ", c.LocalStorage.PasswordStorage)
	log.Println("Data Storage: ", c.LocalStorage.DataStorage)
	log.Println("Text Storage: ", c.LocalStorage.TextStorage)
}

func (c *Client) UnmarshalProtoData(val *pb.CipheredData) (interface{}, error) {
	switch val.Type.String() {
	case "PASSWORD":
		data := models.Password{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case "TEXT":
		data := models.Text{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case "CC":
		data := models.CreditCard{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case "DATA":
		data := models.Data{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("type unknown")
}

func (c *Client) UserRegister(ctx context.Context, user models.User) error {
	err := user.HashPassword()
	if err != nil {
		return err
	}
	userProto := user.ToProto()
	_, err = c.authClient.UserRegister(ctx, &pb.UserRegisterRequest{User: userProto})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UserLogin(ctx context.Context, user models.User) error {
	userProto := user.ToProto()
	response, err := c.authClient.UserLogin(ctx, &pb.UserLoginRequest{User: userProto})
	if err != nil {
		return err
	}
	c.auth.Token = response.Token.Token
	c.auth.Email = response.Token.Email
	c.auth.Secret = response.User.Secret
	go c.RefreshToken(ctx)
	c.crypto = *crypto.NewCrypto(c.auth.Secret)
	user.FromProto(response.User)
	c.currentUser = user
	return nil
}

func (c *Client) RefreshToken(ctx context.Context) {
	tickerRefresh := time.NewTicker(time.Second * 45)
	defer tickerRefresh.Stop()
	for {
		select {
		case <-tickerRefresh.C:
			tokenOld := &pb.Token{
				Email: c.auth.Email,
				Token: c.auth.Token,
			}
			newToken, err := c.authClient.Refresh(ctx, &pb.RefreshRequest{Token: tokenOld})
			if err != nil {
				log.Println("Token not refreshed: ", err)
			}
			c.auth.Token = newToken.Token.Token
		}
	}
}

func (c *Client) GetAllDataFromDB(ctx context.Context) error {
	userProto := c.currentUser.ToProto()
	user := &models.CipheredData{}
	user.User = userProto.Email
	userData := user.ToProto()
	jData, err := c.serverClient.GetCipheredDataForUserRequest(ctx, &pb.GetCipheredDataRequest{Data: userData})
	if err != nil {
		return err
	}
	for _, val := range jData.Data {
		data, err := c.UnmarshalProtoData(val)
		if err != nil {
			return err
		}
		c.AddDataToLocalStorage(ctx, data)
	}
	return nil
}

func (c *Client) AddDataToLocalStorage(ctx context.Context, v any) {
	c.LocalStorage.AppendOrUpdate(v)

}

func (c *Client) AddData(ctx context.Context, v any) error {
	switch v := v.(type) {
	case models.Data:
		if v.ID == "" {
			v.ID = uuid.NewString()
		}
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}

		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "DATA", v.ID)
		_, err = c.serverClient.AddCipheredData(ctx, &pb.AddCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}

		c.AddDataToLocalStorage(ctx, v)
		return nil
	case models.Text:
		if v.ID == "" {
			v.ID = uuid.NewString()
		}
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}

		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "TEXT", v.ID)
		_, err = c.serverClient.AddCipheredData(ctx, &pb.AddCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}
		c.AddDataToLocalStorage(ctx, v)
		return nil
	case models.CreditCard:
		if v.ID == "" {
			v.ID = uuid.NewString()
		}
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}
		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "CC", v.ID)
		_, err = c.serverClient.AddCipheredData(ctx, &pb.AddCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}
		c.AddDataToLocalStorage(ctx, v)
		return nil
	case models.Password:
		if v.ID == "" {
			v.ID = uuid.NewString()
		}
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}
		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "PASSWORD", v.ID)
		_, err = c.serverClient.AddCipheredData(ctx, &pb.AddCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}
		c.AddDataToLocalStorage(ctx, v)
		return nil
	}
	return errors.New("unknown type")
}

func (c *Client) DelData(ctx context.Context, v any) error {
	switch v := v.(type) {
	case models.Data:
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}
		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "DATA", v.ID)
		_, err = c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil
	case models.Password:
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}
		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "PASSWORD", v.ID)
		_, err = c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil
	case models.CreditCard:
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}
		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "CC", v.ID)
		_, err = c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil
	case models.Text:
		dataJson, err := json.Marshal(&v)
		dataJson = c.crypto.Encrypt(dataJson)
		if err != nil {
			return err
		}
		protoData := models.NewCipheredData(dataJson, c.currentUser.Email, "TEXT", v.ID)
		_, err = c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Data: protoData})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil

	}

	return nil
}
