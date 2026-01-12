package grpcclients

import (
	"os"

	pbAuth "yaak-kaii/shared/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client pbAuth.UserServiceClient
	Con    *grpc.ClientConn
}

type GrpcClients struct {
	Auth AuthClient
}

func NewGrpcClient() (*GrpcClients, error) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	if authServiceUrl == "" {
		authServiceUrl = "auth-service:9090"
	}

	connAuth, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	authClient := pbAuth.NewUserServiceClient(connAuth)

	return &GrpcClients{
		Auth: AuthClient{
			Client: authClient,
			Con:    connAuth,
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
