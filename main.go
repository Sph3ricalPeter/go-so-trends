package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	_ "github.com/joho/godotenv/autoload"
)

var (
	neo4jPassword = os.Getenv("NEO4J_PASSWORD")
)

func main() {
	dbUri := "neo4j://127.0.0.1:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", neo4jPassword, ""))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	defer driver.Close(ctx)
	item, err := insertItem(ctx, driver)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", item)
}

func insertItem(ctx context.Context, driver neo4j.DriverWithContext) (*Item, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		"CREATE (n:Item { id: $id, name: $name }) RETURN n",
		map[string]any{
			"id":   1,
			"name": "Item 1",
		}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](result.Records[0], "n")
	if err != nil {
		return nil, fmt.Errorf("could not find node n")
	}
	id, err := neo4j.GetProperty[int64](itemNode, "id")
	if err != nil {
		return nil, err
	}
	name, err := neo4j.GetProperty[string](itemNode, "name")
	if err != nil {
		return nil, err
	}
	return &Item{Id: id, Name: name}, nil
}

type Item struct {
	Id   int64
	Name string
}

func (i *Item) String() string {
	return fmt.Sprintf("Item (id: %d, name: %q)", i.Id, i.Name)
}
