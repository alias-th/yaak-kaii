package main

import (
	"log"
	"yaak-kaii/services/api-gateway/internal/api"
	"yaak-kaii/services/api-gateway/internal/auth"
	grpcclients "yaak-kaii/services/api-gateway/internal/grpc_clients"
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

	app := &api.Application{
		Config: api.Config{
			Addr: httpAddr,
		},
		GrpcClients: grpcClients,
		JwtAuth:     jwtAuth,
	}

	app.Run()

}
