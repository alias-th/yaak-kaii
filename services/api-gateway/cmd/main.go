package main

import (
	"log"
	"yaak-kaii/services/api-gateway/internal/api"
	grpcclients "yaak-kaii/services/api-gateway/internal/grpc_clients"
	"yaak-kaii/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {

	// init grpc clients
	grpcClients, err := grpcclients.NewGrpcCLient()
	if err != nil {
		log.Fatal(err)
	}
	defer grpcClients.Close()

	app := &api.Application{
		Config: api.Config{
			Addr: httpAddr,
		},
		GrpcClients: grpcClients,
	}

	app.Run()

}
