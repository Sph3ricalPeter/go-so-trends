package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Sph3ricalPeter/go-so-trends/api"
	"github.com/Sph3ricalPeter/go-so-trends/internal/config"
	n4j "github.com/Sph3ricalPeter/go-so-trends/internal/db/neo4j"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// connect to DB
	db := &n4j.Neo4j{
		Host:     config.DB_HOST,
		Port:     config.DB_PORT,
		Password: config.DB_PASS,
	}
	driver, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)

	// create repository & api
	repo := &n4j.SoNodeRepository{Driver: driver}
	api := api.SoTrendsApi{Repository: repo, Context: ctx}

	http.Handle("/", api.MountRoutes())
	fmt.Printf("Server is running on %s:%s\n", config.HOST, config.PORT)
	http.ListenAndServe(fmt.Sprintf("%s:%s", config.HOST, config.PORT), nil)
}
