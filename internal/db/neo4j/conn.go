package neo4j

import (
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4j struct {
	Host     string
	Port     string
	Password string
}

func (c *Neo4j) Connect(ctx context.Context) (neo4j.DriverWithContext, error) {
	dbUri := fmt.Sprintf("bolt://%s:%s", c.Host, c.Port)
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", c.Password, ""))
	if err != nil {
		return nil, err
	}
	log.Printf("Connecting to Neo4j at %s\n", dbUri)
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Neo4j")
	return driver, nil
}
