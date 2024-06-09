package neo4j

import (
	"context"
	"fmt"
	"log"

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
	CreateLink(ctx context.Context, source, target string, value float64) error

	// DeleteAll deletes all nodes and links
	DeleteAll(ctx context.Context) error

	// FindTopTags finds the top 5 most popular tags based on node size
	FindTopTags(ctx context.Context, limit int) ([]Node, error)

	// FindTopTagsByDegreeCentrality finds the top 5 most popular tags based on degree centrality using link weights
	FindTopTagsByDegreeCentrality(ctx context.Context, limit int) ([]Node, error)

	// RecommendTags recommends tags based on the given tag and vagueness level
	RecommendTags(ctx context.Context, tag string, vagueness, limit int) ([]Node, error)

	// RecommendTagsSimple recommends tags based on the given tag
	RecommendTagsSimple(ctx context.Context, tag string, limit int) ([]string, error)
}

func (s *SoNodeRepository) NewSoNodeRepository(driver neo4j.DriverWithContext) *SoNodeRepository {
	return &SoNodeRepository{Driver: driver}
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

func (s *SoNodeRepository) CreateLink(ctx context.Context, source, target string, value float64) error {
	_, err := neo4j.ExecuteQuery(ctx, s.Driver,
		"MATCH (source:Node { name: $source }), (target:Node { name: $target }) CREATE (source)-[:LINK { value: $value }]->(target)",
		map[string]interface{}{
			"source": source,
			"target": target,
			"value":  value,
		}, neo4j.EagerResultTransformer)
	return err
}

func (s *SoNodeRepository) DeleteAll(ctx context.Context) error {
	_, err := neo4j.ExecuteQuery(ctx, s.Driver, "MATCH (n) DETACH DELETE n", nil, neo4j.EagerResultTransformer)
	return err
}

func (s *SoNodeRepository) FindTopTags(ctx context.Context, limit int) ([]Node, error) {
	result, err := neo4j.ExecuteQuery(ctx, s.Driver,
		"MATCH (n:Node) RETURN n ORDER BY n.size DESC LIMIT $limit",
		map[string]any{
			"limit": limit,
		},
		neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	return createNodesFromResult(result)
}

func (s *SoNodeRepository) FindTopTagsByDegreeCentrality(ctx context.Context, limit int) ([]Node, error) {
	result, err := neo4j.ExecuteQuery(ctx, s.Driver,
		"MATCH (n:Node)-[r:LINK]->() RETURN n, SUM(r.weight) AS degree ORDER BY degree DESC LIMIT $limit",
		map[string]any{
			"limit": limit,
		},
		neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	return createNodesFromResult(result)
}

func (s *SoNodeRepository) RecommendTags(ctx context.Context, tag string, vagueness, limit int) ([]Node, error) {
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := fmt.Sprintf(
		`MATCH path = (n:Node {name: "%s"})-[*1..%d]->(m:Node)
		WHERE n <> m
		WITH m, path, reduce(valueSum = 0, r in relationships(path) | valueSum + coalesce(r.value, 0)) AS totalValue
		WHERE totalValue >= %.2f
		RETURN m, totalValue
		ORDER BY totalValue DESC
		LIMIT %d`, tag, vagueness, 0.0, limit)

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0)
	for result.Next(ctx) {
		record := result.Record()
		log.Printf("record: %v", record.Keys)
		nodeItem, ok := record.Get("m")
		if !ok {
			return nil, fmt.Errorf("could not find node n")
		}
		item, ok := nodeItem.(neo4j.Node)
		if !ok {
			return nil, fmt.Errorf("could not find node n")
		}
		props := item.Props
		name, ok := props["name"]
		if !ok {
			return nil, fmt.Errorf("could not find node name")
		}
		group, ok := props["group"]
		if !ok {
			return nil, fmt.Errorf("could not find node group")
		}
		size, ok := props["size"]
		if !ok {
			return nil, fmt.Errorf("could not find node size")
		}
		metadata := make(map[string]interface{})
		degree, ok := record.Get("totalValue")
		if ok {
			metadata["degree"] = degree
		}
		nodes = append(nodes, Node{Name: name.(string), Group: group.(int64), Size: size.(float64), Metadata: metadata})
	}
	return nodes, nil
}

func (s *SoNodeRepository) RecommendTagsSimple(ctx context.Context, tag string, limit int) ([]string, error) {
	session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := fmt.Sprintf(
		`MATCH path = (n:Node {name: "%s"})-[*1]->(m:Node)
		WHERE n <> m
		WITH m, path, reduce(valueSum = 0, r in relationships(path) | valueSum + coalesce(r.value, 0)) AS totalValue
		RETURN m, totalValue
		ORDER BY totalValue DESC
		LIMIT %d`, tag, limit)

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	for result.Next(ctx) {
		record := result.Record()
		log.Printf("record: %v", record.Keys)
		nodeItem, ok := record.Get("m")
		if !ok {
			return nil, fmt.Errorf("could not find node n")
		}
		item, ok := nodeItem.(neo4j.Node)
		if !ok {
			return nil, fmt.Errorf("could not find node n")
		}
		props := item.Props
		name, ok := props["name"]
		if !ok {
			return nil, fmt.Errorf("could not find node name")
		}
		tags = append(tags, name.(string))
	}
	return tags, nil
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
		metadata := make(map[string]interface{})
		degree, _, err := neo4j.GetRecordValue[int64](record, "degree")
		if err == nil {
			metadata["degree"] = degree
		}
		nodes = append(nodes, Node{Name: name, Group: group, Size: size, Metadata: metadata})
	}

	return nodes, nil
}

func (n *Node) String() string {
	return fmt.Sprintf("Node (Name: %s, Group: %d, Size: %f)", n.Name, n.Group, n.Size)
}

type SoNodeRepository struct {
	Driver neo4j.DriverWithContext
}

type Node struct {
	Name     string                 `json:"name"`
	Group    int64                  `json:"group"`
	Size     float64                `json:"size"`
	Metadata map[string]interface{} `json:"metadata"`
}
