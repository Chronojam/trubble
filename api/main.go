package main

import (
	"golang.org/x/net/context"
	"log"
	"net"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"

	"github.com/chronojam/trubble/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/mwitkow/go-grpc-middleware"
	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/chronojam/trubble/api/proto"
)

var dsClient *datastore.Client
var psClient *pubsub.Client
var jwtSecret = "HelloWorld"

func main() {
	ctx := context.Background()
	// Setup datastore connection
	d, err := datastore.NewClient(ctx, "chronojam-trubble")
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Setup pubsub Client
	p, err := pubsub.NewClient(ctx, "chronojam-trubble")
	if err != nil {
		log.Fatalf(err.Error())
	}

	dsClient = d
	psClient = p

	conn, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatalf(err.Error())
	}

	connUser, err := net.Listen("tcp", ":3334")
	if err != nil {
		log.Fatalf(err.Error())
	}

	connAdmin, err := net.Listen("tcp", ":3335")
	if err != nil {
		log.Fatalf(err.Error())
	}

	noTokenServer := grpc.NewServer()
	adminServer := grpc.NewServer(grpc.UnaryInterceptor(AuthUnaryInterceptorAdmin))
	tokenServer := grpc.NewServer(grpc.StreamInterceptor(AuthStreamInterceptorUser), grpc.UnaryInterceptor(AuthUnaryInterceptorUser))

	pb.RegisterAuthServer(noTokenServer, &AuthServer{})
	pb.RegisterAdminServer(adminServer, &AdminServer{})
	pb.RegisterGameServer(tokenServer, &GameServer{})

	go noTokenServer.Serve(conn)
	go tokenServer.Serve(connUser)
	adminServer.Serve(connAdmin)
}

// AuthStreamInterceptorUser ...
// Handles token authentication
func AuthStreamInterceptorUser(
	req interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	ctx := stream.Context()
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return errors.New("error parsing grpc metadata")
	}

	if _, ok := md["jwt"]; !ok {
		return errors.New("authorization header missing")
	}
	tok, err := util.ParseJwtToken(md["jwt"][0], jwtSecret)
	if err != nil {
		return errors.Wrap(err, "bad token")
	}
	// add values to context here.
	claims := tok.Claims.(jwt.MapClaims)
	accountID := claims["_account_id"].(float64)

	newCtx := context.WithValue(ctx, "_account_id", accountID)

	wrapped := grpc_middleware.WrapServerStream(stream)
	wrapped.WrappedContext = newCtx
	// scope auth here
	return handler(req, wrapped)
}

// AuthUnaryInterceptorUser ...
// Handles token authentication
func AuthUnaryInterceptorUser(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return req, errors.New("error parsing grpc metadata")
	}

	if _, ok := md["jwt"]; !ok {
		return req, errors.New("authorization header missing")
	}
	tok, err := util.ParseJwtToken(md["jwt"][0], jwtSecret)
	if err != nil {
		return req, errors.Wrap(err, "bad token")
	}
	// add values to context here.
	claims := tok.Claims.(jwt.MapClaims)
	accountID := claims["_account_id"].(float64)

	newCtx := context.WithValue(ctx, "_account_id", accountID)

	// scope auth here
	return handler(newCtx, req)
}

// AuthUnaryInterceptorAdmin ...
// Handles token authentication for admins
func AuthUnaryInterceptorAdmin(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return req, errors.New("error parsing grpc metadata")
	}

	if _, ok := md["jwt"]; !ok {
		return req, errors.New("authorization header missing")
	}
	tok, err := util.ParseJwtToken(md["jwt"][0], jwtSecret)
	if err != nil {
		return req, errors.Wrap(err, "bad token")
	}
	// add values to context here.
	claims := tok.Claims.(jwt.MapClaims)
	accountID := claims["_account_id"].(float64)

	// only admins can access these routes
	// check if we have this claim
	if _, ok := claims["_admin"]; !ok {
		return req, errors.New("unauthorized")
	}

	// check that the claim is a boolean
	if _, ok := claims["_admin"].(bool); !ok {
		return req, errors.New("unauthorized")
	}

	admin := claims["_admin"].(bool)

	// check that the value is true.
	if !admin {
		return req, errors.New("unauthorized")
	}

	newCtx := context.WithValue(ctx, "_account_id", accountID)
	newCtx = context.WithValue(newCtx, "_admin", admin)

	// scope auth here
	return handler(newCtx, req)
}
