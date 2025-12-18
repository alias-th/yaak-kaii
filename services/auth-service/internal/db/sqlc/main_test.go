package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	connString := "postgres://root:password123@localhost:5432/yaak_kaii?sslmode=disable"

	connPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	log.Println("connected to database")
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
