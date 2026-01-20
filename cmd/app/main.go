package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	httpTransport.RegisterRoutes(engine, repo, genID)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", engine); err != nil {
		log.Fatal(err)
	}
}
