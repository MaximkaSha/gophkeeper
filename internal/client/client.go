// Package client implements gRPC client.
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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// UnaryAuthClientInterceptor - auth middleware. Adds authorization header to each client request.
func (a *Auth) UnaryAuthClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+a.Token)
	return invoker(ctx, method, req, reply, cc, opts...)
}

// Auth structure used to keep users auth info.
type Auth struct {
	// JWT Token.
	Token string
	// User email.
	Email string
	// bcrypt hash of user's password.
	Hash string
	// Personal user key to crypt data.
	Secret []byte
}

// LocalStorage struct used to keep models data for UI.
// This structed used to keep work with UI tables simple.
type LocalStorage struct {
	DataStorage     []models.Data
	TextStorage     []models.Text
	CCStorage       []models.CreditCard
	PasswordStorage []models.Password
}

// AllData structure keeps all user data in one place.
// This is used to keep user data returned from server.
type AllData struct {
	ID    string
	JData []byte
	Type  string
}

// AppendOrUpdate func append data to localstorage for UI.
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

// DelFromLocalStorageUI - Delete data from LocalStorage.
func (l *LocalStorage) DelFromLocalStorageUI(v any) {
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

// NewLocalStorage - LocalStorage constructor.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		// Slice of files model.
		DataStorage: []models.Data{},
		// Slice of texts model.
		TextStorage: []models.Text{},
		// Slice of cc model.
		CCStorage: []models.CreditCard{},
		// Slice of password model.
		PasswordStorage: []models.Password{},
	}
}

// Client struct main structure of client app.
type Client struct {
	authClient   pb.AuthGophkeeperClient
	serverClient pb.GophkeeperClient
	crypto       crypto.Crypto
	currentUser  models.User
	auth         *Auth
	// LocalStorage for UI.
	LocalStorage *LocalStorage
	// Configuration data.
	Config *config.ClientConfig
	// Build version from linker.
	BuildVersion string
	// Build time from linker.
	BuildTime string
	// All data slice.
	AllData []AllData
}

// NewClient Client constructor.
// Build version and time must be passed as string.
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
		AllData:      []AllData{},
	}
}

// PrinStorage print to log current state of local storage.
// Debug purpose.
func (c *Client) PrinStorage() {
	log.Println("CC Storage: ", c.LocalStorage.CCStorage)
	log.Println("Password Storage: ", c.LocalStorage.PasswordStorage)
	log.Println("Data Storage: ", c.LocalStorage.DataStorage)
	log.Println("Text Storage: ", c.LocalStorage.TextStorage)
}

// UnmarshalProtoData function decrypt and unmarshal data from protobuf to models.
func (c *Client) UnmarshalProtoData(val *pb.CipheredData) (interface{}, error) {
	switch val.Type.String() {
	case "PASSWORD":
		data := models.Password{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		data.ID = val.Uuid
		return data, nil
	case "TEXT":
		data := models.Text{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		data.ID = val.Uuid
		return data, nil
	case "CC":
		data := models.CreditCard{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		data.ID = val.Uuid
		return data, nil
	case "DATA":
		data := models.Data{}
		val.Data = c.crypto.Decrypt(val.Data)
		err := json.Unmarshal(val.Data, &data)
		if err != nil {
			return nil, err
		}
		data.ID = val.Uuid
		return data, nil
	}
	return nil, errors.New("type unknown")
}

// UserRegister - registration function.
// models.User must be passed.
// Return error if error occures when writing to DB (eg. User already exist).
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

// UserLogin - login function.
// models.User must be passed.
// Return error if error occures when writing to DB (eg. Bad pwd).
// If all ok  jwt token and privite key will placed to Client object.
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

// RefreshToken - refresh token every 45 seconds.
func (c *Client) RefreshToken(ctx context.Context) {
	tickerRefresh := time.NewTicker(time.Second * 45)
	defer tickerRefresh.Stop()
	for range tickerRefresh.C {
		tokenOld := &pb.Token{
			Email: c.auth.Email,
			Token: c.auth.Token,
		}
		newToken, err := c.authClient.Refresh(ctx, &pb.RefreshRequest{Token: tokenOld})
		if err != nil {
			log.Println("Token not refreshed: ", err)
			return
		}
		c.auth.Token = newToken.Token.Token
	}
}

// GetAllDataFromDB - ask server for all users data in DB.
// Data will be writen to AllData slice.
func (c *Client) GetAllDataFromDB(ctx context.Context) error {
	jData, err := c.serverClient.GetCipheredDataForUserRequest(ctx, &pb.GetCipheredDataRequest{Email: c.currentUser.Email})
	if err != nil {
		return err
	}
	for _, val := range jData.Data {
		data, err := c.UnmarshalProtoData(val)
		if err != nil {
			return err
		}
		c.AddDataToLocalStorage(ctx, data.(models.Dater))
	}
	return nil
}

// AddDataToLocalStorageUI - Adds data to local storage for UI
func (c *Client) AddDataToLocalStorageUI(ctx context.Context, v any) {
	c.LocalStorage.AppendOrUpdate(v)

}

// AddDataToLocalStorage - Adds data to AllData storage.
func (c *Client) AddDataToLocalStorage(ctx context.Context, data models.Dater) {

	for i := range c.AllData {
		if c.AllData[i].ID == data.GetID() {
			c.AllData[i].JData = data.GetData()
			return
		}
	}
	c.AllData = append(c.AllData, AllData{
		ID:    data.GetID(),
		JData: data.GetData(),
		Type:  data.Type(),
	})
}

// DelFromLocalStorage -Del data from AllData storage by given uuid.
func (c *Client) DelFromLocalStorage(uuid string) {
	if len(c.AllData) == 1 {
		c.AllData = make([]AllData, 0)
	}
	for i := range c.AllData {
		if c.AllData[i].ID == uuid {
			c.AllData = append(c.AllData[:i], c.AllData[i+1:]...)
			return
		}
	}
}

// AddData - encrypt  and push data to server.
func (c *Client) AddData(ctx context.Context, data models.Dater) error {
	cData := c.crypto.Encrypt(data.GetData())
	protoData := models.NewCipheredData(cData, c.currentUser.Email, data.Type(), data.GetID())
	_, err := c.serverClient.AddCipheredData(ctx, &pb.AddCipheredDataRequest{Data: protoData})
	if err != nil {
		return err
	}
	c.AddDataToLocalStorage(ctx, data)
	return nil
}

/*func (c *Client) AddData_Old(ctx context.Context, v any) error {
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
} */

// DelData - delete data from server and local storage by given uuid.
func (c *Client) DelData(ctx context.Context, uuid string) error {
	_, err := c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Uuid: uuid})
	if err != nil {
		return err
	}
	c.DelFromLocalStorage(uuid)
	return nil
}

/*
func (c *Client) DelData_Old(ctx context.Context, v any) error {
	switch v := v.(type) {
	case models.Data:
		_, err := c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Uuid: v.ID})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil
	case models.Password:
		_, err := c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Uuid: v.ID})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil
	case models.CreditCard:
		_, err := c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Uuid: v.ID})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil
	case models.Text:
		_, err := c.serverClient.DelCipheredData(ctx, &pb.DelCipheredDataRequest{Uuid: v.ID})
		if err != nil {
			return err
		}
		c.LocalStorage.DelFromLocalStorage(v)
		return nil

	}

	return nil
} */
