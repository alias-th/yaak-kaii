package main

import (
	"log"
	"yaak-kaii/services/api-gateway/internal/api"
	"yaak-kaii/services/api-gateway/internal/auth"
	grpcclients "yaak-kaii/services/api-gateway/internal/grpc_clients"
	"yaak-kaii/services/api-gateway/internal/s3"
	"yaak-kaii/shared/env"
)

var (
	httpAddr  = env.GetString("HTTP_ADDR", ":8081")
	jwtSecret = env.GetString("JWT_SECRET", "SECRET")
	jwtISS    = env.GetString("JWT_ISSUER", "MONTON")
)

func main() {

	// init grpc clients
	grpcClients, err := grpcclients.NewGrpcClient()
	if err != nil {
		log.Fatal(err)
	}
	defer grpcClients.Close()

	jwtAuth := auth.NewJWTAuthenticator(jwtSecret, jwtISS, jwtISS)

	s3AccountID := env.GetString("S3_ACCOUNT_ID", "")
	s3AccessKey := env.GetString("S3_ACCESS_KEY_ID", "")
	s3SecretKey := env.GetString("S3_SECRET_ACCESS_KEY", "")
	s3Bucket := env.GetString("S3_BUCKET_NAME", "")
	s3Uploader, err := s3.NewS3Uploader(s3AccountID, s3AccessKey, s3SecretKey, s3Bucket)
	if err != nil {
		log.Fatal("failed to initialize S3 uploader:", err)
	}

	app := &api.Application{
		Config: api.Config{
			Addr: httpAddr,
		},
		GrpcClients: grpcClients,
		JwtAuth:     jwtAuth,
		S3Uploader:  s3Uploader,
	}

	app.Run()

}
