package config

import (
	"os"
	"yaak-kaii/shared/env"
)

type config struct {
	ConnString string
	GrpcAddr   string
	JwtSecret  string
	JwtISS     string
}

func LoadConfig() config {
	return config{
		ConnString: os.Getenv("PG_URI"),
		GrpcAddr:   ":9090",
		JwtSecret:  env.GetString("JWT_SECRET", "SECRET"),
		JwtISS:     env.GetString("JWT_ISSUER", "MONTON"),
	}
}
