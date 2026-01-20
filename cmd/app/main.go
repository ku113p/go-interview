package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-interview/internal/biography/infra/postgres"
	httpTransport "go-interview/internal/biography/transport/http"
	"go-interview/pkg/utils"
)

func main() {
	db, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := postgres.NewAreaRepository(db)
	genID := utils.NewUUID7Generator()

	router := httpTransport.NewRouter(repo, genID)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
