package main

import (
	"golang.org/x/net/context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"cloud.google.com/go/datastore"

	"github.com/pkg/errors"
	"github.com/dgrijalva/jwt-go"
	proto "github.com/chronojam/trubble/api/proto"
)

// AuthServer ...
// GRPC server object.
type AuthServer struct {}

// IssueNewJWT ...
// Issues a new JWT token on validation of user credentials
func (a *AuthServer) IssueNewJWT(ctx context.Context, req *proto.NewJwtRequest) (*proto.NewJwtResponse, error) {
	// Validation
	if req.GetEmail() == "" {
		return nil, errors.New("Email is a required field")
	}

	if req.GetPassword() == "" {
		return nil, errors.New("Password is a required field")
	}

	email := req.GetEmail()
	password := req.GetPassword()

	var accounts []Account

	query := datastore.NewQuery("Account").
	Filter("Email =", email).
	Limit(1)

	_, err := dsClient.GetAll(ctx, query, &accounts)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting accounts from datastore")
	}

	if len(accounts) != 1 {
		return nil, errors.New("error while getting accounts from datastore, accounts expected == 1")
	}

	err = bcrypt.CompareHashAndPassword([]byte(accounts[0].Password), []byte(password))
	if err != nil {
		return nil, errors.Wrap(err, "incorrect username or password")
	}

	admin := false
	if (accounts[0].Admin) {
		admin = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"_account_id": accounts[0].K.ID,
		"_admin": admin,
		"exp": time.Now().Add(time.Hour).Unix(),
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.Wrap(err, "could not create jwt token")
	}

	resp := &proto.NewJwtResponse{
		Token: tokenString,
	}

	return resp, nil
}

// UpdatePassword ...
// Issues a new JWT token on validation of user credentials
func (a *AuthServer) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordRequest) (*proto.UpdatePasswordResponse, error) {
	// TODO Implement this.
	return nil, nil
}