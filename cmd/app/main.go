package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	create_life_area "go-interview/internal/biography/app/commands/create_life_area"
	get_life_area "go-interview/internal/biography/app/queries/get_life_area"
	"go-interview/internal/biography/infra/postgres"
	"go-interview/pkg/utils"
)

func main() {
	db, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	areaRepo := postgres.NewAreaRepository(db)
	genID := utils.NewUUID7Generator()

	createHandler := create_life_area.NewCreateLifeAreaHandler(areaRepo, genID)
	getHandler := get_life_area.NewGetLifeAreaHandler(areaRepo)

	createCmd := create_life_area.CreateLifeAreaCommand{
		UserID:   uuid.New().String(),
		ParentID: nil,
		Title:    "New Life Area from command",
		Goal:     "Initial Goal from command",
	}

	createResult, err := createHandler.Handle(context.Background(), createCmd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created life area: %s\n", createResult)

	getQuery := get_life_area.GetLifeAreaQuery{
		ID: createResult.ID,
	}

	getResult, err := getHandler.Handle(context.Background(), getQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("retrieved life area: %+v\n", getResult)
}
