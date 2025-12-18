package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testStore Store

func TestMain(m *testing.M) {
	godotenv.Load("../../../.env.local")
	connString := os.Getenv("PG_URI")

	connPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	log.Println("connected to database")
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
