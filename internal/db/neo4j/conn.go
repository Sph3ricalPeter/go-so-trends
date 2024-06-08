package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Db interface {
	Connect() (neo4j.DriverWithContext, error)
}

type Neo4j struct {
	Host     string
	Port     string
	Password string
}

func (c *Neo4j) Connect() (neo4j.DriverWithContext, error) {
	dbUri := fmt.Sprintf("neo4j://%s:%s", c.Host, c.Port)
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", c.Password, ""))
	if err != nil {
		return nil, err
	}
	return driver, nil
}
