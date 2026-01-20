package main

import (
	"context"
	"log"
	"net/http"

	httpTransport "go-interview/internal/biography/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"

	create_life_area "go-interview/internal/biography/app/commands/create_life_area"
	"go-interview/internal/biography/infra/postgres"
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

	useCase := create_life_area.NewCreateLifeAreaHandler(repo, genID)
	httpHandler := httpTransport.NewCreateLifeAreaHandlerHTTP(useCase)
	http.HandleFunc("/life-areas", httpHandler.Handle)

	http.ListenAndServe(":8080", nil)
}
