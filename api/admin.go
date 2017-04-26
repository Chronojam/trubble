package main

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/chronojam/trubble/util"
	proto "github.com/chronojam/trubble/api/proto"
)

// AdminServer ...
// GRPC Server object for administrative functions
type AdminServer struct {}

// CreateAccount ...
// Creates a new account with the given GRPC Request
func (a *AdminServer) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	// Validation
	if req.GetEmail() == "" {
		return nil, errors.New("email is a required field")
	}

	email := req.GetEmail()

	// First we'll check if there is another user with that email
	// We'd like these to be unique, so we'll check both
	var accounts []Account
	query := datastore.NewQuery("Account").
	Filter("Email = ", email).
	Limit(1)

	_, err := dsClient.GetAll(ctx, query, &accounts)
	if err != nil {
		return nil, err
	}

	if len(accounts) != 0 {
		// So this username/email combo already exists.
		// Note that we are looking for an account where the email and username matches.
		return nil, errors.New("account already exists. Forgotten password?")
	}

	// Create a temporary password.
	pass := util.RandomBytes(24)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Store new user in datastore.
	key := datastore.NameKey("Account", "", nil)
	account := Account{
		Email: email,
		Password: string(hash),
	}

	_, err = dsClient.Put(ctx, key, &account)
	if err != nil {
		return nil, err
	}

	resp := &proto.CreateAccountResponse{
		Temppassword: string(pass),
	}

	return resp, nil
} 

// DeleteAccount ...
// Deletes the given account by ID
func (a *AdminServer) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	// TODO Delete account logic here.
	return nil, nil
}