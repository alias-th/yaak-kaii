package grpcclients

import (
	"os"

	pb "yaak-kaii/shared/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client pb.UserServiceClient
	Con    *grpc.ClientConn
}

type GrpcClients struct {
	Auth AuthClient
}

func NewGrpcCLient() (*GrpcClients, error) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	if authServiceUrl == "" {
		authServiceUrl = "auth-service:9090"
	}

	conn, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	authClient := pb.NewUserServiceClient(conn)

	return &GrpcClients{
		Auth: AuthClient{
			Client: authClient,
			Con:    conn,
		},
	}, nil
}

func (client *GrpcClients) Close() {
	connAuth := client.Auth.Con
	if connAuth != nil {
		if err := connAuth.Close(); err != nil {
			return
		}
	}
}
