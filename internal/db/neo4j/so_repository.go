package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type NodeRepository interface {
	// InsertNode inserts a new node into the database, returns the inserted node
	InsertNode(ctx context.Context, node Node) (*Node, error)

	// FindNodesByName finds all nodes with given name, returns error if not found
	FindNodesByName(ctx context.Context, name string) ([]Node, error)

	// FindNodeByName finds first node with given name, returns error if not found
	FindNodeByName(ctx context.Context, name string) (*Node, error)

	// CreateLink creates a link between two nodes
	CreateLink(ctx context.Context, source, target string) error

	// DeleteAll deletes all nodes and links
	DeleteAll(ctx context.Context) error

	// FindTopTags finds the top 5 most popular tags based on node size
	FindTopTags(ctx context.Context) ([]Node, error)
}

type SoNodeRepository struct {
	Driver neo4j.DriverWithContext
}

type Node struct {
	Name  string  `json:"name"`
	Group int64   `json:"group"`
	Size  float64 `json:"size"`
}

func (s *SoNodeRepository) InsertNode(ctx context.Context, node Node) (*Node, error) {
	result, err := neo4j.ExecuteQuery(ctx, s.Driver,
		"CREATE (n:Node { name: $name, group: $group, size: $size }) RETURN n",
		map[string]any{
			"name":  node.Name,
			"group": node.Group,
			"size":  node.Size,
		}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	nodes, err := createNodesFromResult(result)
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

func (s *SoNodeRepository) FindNodeByName(ctx context.Context, name string) (*Node, error) {
	nodes, err := s.FindNodesByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("node with name %s not found", name)
	}
	return &nodes[0], nil
}

func (s *SoNodeRepository) FindNodesByName(ctx context.Context, name string) ([]Node, error) {
	result, err := neo4j.ExecuteQuery(ctx, s.Driver,
		"MATCH (n) WHERE n.name = $name RETURN n",
		map[string]any{
			"name": name,
		}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	return createNodesFromResult(result)
}

func (s *SoNodeRepository) CreateLink(ctx context.Context, source, target string) error {
	_, err := neo4j.ExecuteQuery(ctx, s.Driver,
		"MATCH (source:Node { name: $source }), (target:Node { name: $target }) CREATE (source)-[:LINK]->(target)",
		map[string]any{
			"source": source,
			"target": target,
		}, neo4j.EagerResultTransformer)
	return err
}

func (s *SoNodeRepository) DeleteAll(ctx context.Context) error {
	_, err := neo4j.ExecuteQuery(ctx, s.Driver, "MATCH (n) DETACH DELETE n", nil, neo4j.EagerResultTransformer)
	return err
}

func (s *SoNodeRepository) FindTopTags(ctx context.Context) ([]Node, error) {
	result, err := neo4j.ExecuteQuery(ctx, s.Driver,
		"MATCH (n:Node) RETURN n ORDER BY n.size DESC LIMIT 5",
		nil, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	return createNodesFromResult(result)
}

func createNodesFromResult(result *neo4j.EagerResult) ([]Node, error) {
	nodes := make([]Node, 0)

	for _, record := range result.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "n")
		if err != nil {
			return nil, fmt.Errorf("could not find node n")
		}
		name, err := neo4j.GetProperty[string](itemNode, "name")
		if err != nil {
			return nil, err
		}
		group, err := neo4j.GetProperty[int64](itemNode, "group")
		if err != nil {
			return nil, err
		}
		size, err := neo4j.GetProperty[float64](itemNode, "size")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, Node{Name: name, Group: group, Size: size})
	}

	return nodes, nil
}

func (n *Node) String() string {
	return fmt.Sprintf("Node (Name: %s, Group: %d, Size: %f)", n.Name, n.Group, n.Size)
}
