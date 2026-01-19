package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-interview/internal/biography/application"
	"go-interview/internal/biography/application/dto"
	"go-interview/internal/biography/infrastructure/persistence/postgres"
)

func main() {
	db, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	areaRepo := postgres.NewAreaRepository(db)
	areaService := application.NewAreaService(areaRepo)

	userID := uuid.New()
	newArea, err := areaService.Create(context.Background(), dto.CreateAreaRequest{
		UserID:   userID,
		ParentID: nil,
		Title:    "New Life Area",
		Goal:     "Initial Goal",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created area: %v\n", newArea)

	err = areaService.ChangeGoal(context.Background(), dto.ChangeGoalRequest{
		ID:   newArea.ID,
		Goal: "Updated Goal",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("changed goal of the area")

	areas, err := areaService.List(context.Background(), userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("areas for user %s: %v\n", userID, areas)
}
