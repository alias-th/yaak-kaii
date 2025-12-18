package grpcclients

import (
	"os"

	pb "yaak-kaii/shared/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authServiceClient struct {
	Client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewAuthServiceClient() (*authServiceClient, error) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	if authServiceUrl == "" {
		authServiceUrl = "auth-service:9090"
	}

	conn, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)

	return &authServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *authServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
