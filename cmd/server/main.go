package main

import (
	"context"
	"log"
	"time"

	"github.com/britinogn/quillhub/internal/database"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := database.ConnectPostgres(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	log.Println("Postgres connected successfully with pgx")
}