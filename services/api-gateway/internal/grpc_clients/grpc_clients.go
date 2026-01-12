package grpcclients

import (
	"os"

	pbAuth "yaak-kaii/shared/proto/auth"
	pbProduct "yaak-kaii/shared/proto/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client pbAuth.UserServiceClient
	Con    *grpc.ClientConn
}

type ProductClient struct {
	Client pbProduct.ProductServiceClient
	Con    *grpc.ClientConn
}

type GrpcClients struct {
	Auth    AuthClient
	Product ProductClient
}

func NewGrpcClient() (*GrpcClients, error) {
	authServiceUrl := os.Getenv("AUTH_SERVICE_URL")
	productServiceUrl := os.Getenv("PRODUCT_SERVICE_URL")
	if authServiceUrl == "" {
		authServiceUrl = "auth-service:9090"
	}
	if productServiceUrl == "" {
		productServiceUrl = "product-service:9091"
	}

	connAuth, err := grpc.NewClient(authServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	connProduct, err := grpc.NewClient(productServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	authClient := pbAuth.NewUserServiceClient(connAuth)
	productClient := pbProduct.NewProductServiceClient(connProduct)

	return &GrpcClients{
		Auth: AuthClient{
			Client: authClient,
			Con:    connAuth,
		},
		Product: ProductClient{
			Client: productClient,
			Con:    connProduct,
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
