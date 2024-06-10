package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Sph3ricalPeter/go-so-trends/internal/config"
	n4j "github.com/Sph3ricalPeter/go-so-trends/internal/db/neo4j"
)

const (
	NODES_FNAME = "data/stack_network_nodes.csv"
	LINKS_FNAME = "data/stack_network_links.csv"
)

func main() {
	// connect to DB
	db := &n4j.Neo4j{
		Host:     config.DB_HOST,
		Port:     config.DB_PORT,
		Password: config.DB_PASS,
	}
	ctx := context.Background()
	driver, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer driver.Close(ctx)

	repo := n4j.SoNodeRepository{Driver: driver}

	// cleanup existing data
	repo.DeleteAll(ctx)

	// load and insert nodes
	nodeRecords := loadCSVRecords(NODES_FNAME)
	linkRecords := loadCSVRecords(LINKS_FNAME)

	insertNodes(ctx, repo, nodeRecords)
	createLinks(ctx, repo, linkRecords)
}

func insertNodes(ctx context.Context, repo n4j.SoNodeRepository, records [][]string) {
	for _, record := range records[1:] {
		node, err := parseNodeFromRecord(record)
		if err != nil {
			panic(err)
		}

		_, err = repo.InsertNode(ctx, node)
		if err != nil {
			panic(err)
		}
	}
}

func createLinks(ctx context.Context, repo n4j.SoNodeRepository, records [][]string) {
	for _, record := range records[1:] {
		source, err := repo.FindNodeByName(ctx, record[0])
		if err != nil {
			panic(err)
		}

		target, err := repo.FindNodeByName(ctx, record[1])
		if err != nil {
			panic(err)
		}

		value, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			panic(err)
		}

		err = repo.CreateLink(ctx, source.Name, target.Name, value)
		if err != nil {
			panic(err)
		}
	}
}

func loadCSVRecords(fpath string) [][]string {
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func parseNodeFromRecord(record []string) (n4j.Node, error) {
	group, err := strconv.ParseInt(record[1], 10, 64)
	if err != nil {
		fmt.Printf("error parsing group from record: %v\n", record)
		return n4j.Node{}, err
	}

	size, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return n4j.Node{}, err
	}

	node := n4j.Node{
		Name:  record[0],
		Group: group,
		Size:  size,
	}

	return node, nil
}
