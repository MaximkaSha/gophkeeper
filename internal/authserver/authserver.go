package authserver

import (
	"context"
	"time"

	"github.com/MaximkaSha/gophkeeper/internal/models"
	pb "github.com/MaximkaSha/gophkeeper/internal/proto"
	"github.com/MaximkaSha/gophkeeper/internal/storage"
	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGophkeeperServer struct {
	jwtKey []byte
	DB     *storage.Storage
	pb.UnimplementedAuthGophkeeperServer
}

func NewAuthGophkeeperServer() AuthGophkeeperServer {
	return AuthGophkeeperServer{
		jwtKey: []byte("my_secret_key"), // do I need to add random for each server start ?
		DB:     storage.NewStorage(),
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (a AuthGophkeeperServer) UserRegister(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {

	var response pb.UserRegisterResponse
	user := models.User{}
	user.FromProto(in.User)
	err := a.DB.AddUser(user)
	if err != nil {
		return &response, status.Errorf(codes.AlreadyExists, `User already exists`)
	}
	return &response, nil
}

func (a AuthGophkeeperServer) JWTClain(creds models.User) (string, int64, error) {

	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtKey)
	return tokenString, expirationTime.Unix(), err
}

func (a AuthGophkeeperServer) UserLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	var response pb.UserLoginResponse
	user := models.User{}
	user.FromProto(in.User)
	userPass := models.User{}
	userPass, err := a.DB.GetUser(user)
	if err != nil {
		return &response, status.Errorf(codes.NotFound, err.Error())
	}
	if user.CheckPasswordHash(userPass.Password) {
		return &response, status.Errorf(codes.Unauthenticated, "wrong password")
	}
	tokenString, expiresAt, err := a.JWTClain(userPass)
	if err != nil {
		return &response, status.Errorf(codes.Unknown, "JWT Generating Error")
	}
	token := pb.Token{
		Email:   userPass.Email,
		Token:   tokenString,
		Expires: expiresAt,
	}
	response.User = userPass.ToProto()
	response.Token = &token
	return &response, nil
}

func (a AuthGophkeeperServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (a AuthGophkeeperServer) parseToken(token string) error {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return status.Error(codes.Unauthenticated, "wrong token")
		}
	}
	if !tkn.Valid {
		status.Error(codes.Unauthenticated, "wrong token")
	}

	return nil
}

func (a AuthGophkeeperServer) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	var response pb.RefreshResponse
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(in.Token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return &response, status.Errorf(codes.Unauthenticated, "wrong password")
		}
	}
	if !tkn.Valid {
		return &response, status.Errorf(codes.Unauthenticated, "wrong password")
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return &response, status.Errorf(codes.FailedPrecondition, "too early to refresh")
	}

	expirationTime := time.Now().Add(1 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenResp := in.Token
	tokenResp.Token, err = token.SignedString(a.jwtKey)
	if err != nil {
		return &response, status.Errorf(codes.Unknown, "Internal error")
	}
	response.Token = tokenResp
	return &response, nil
}

// AuthFunc is used by a middleware to authenticate requests.
func (a AuthGophkeeperServer) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	err = a.parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	return ctx, nil
}