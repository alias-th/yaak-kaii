package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	db "yaak-kaii/services/auth-service/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

var (
	connString = os.Getenv("PG_URI")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	_ = db.NewStore(connPool)

}
