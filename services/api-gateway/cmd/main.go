package main

import (
	"yaak-kaii/services/api-gateway/internal/api"
	"yaak-kaii/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	api.RunGinServer(httpAddr)
}
