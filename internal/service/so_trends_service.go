package service

import (
	"context"

	n4j "github.com/Sph3ricalPeter/go-so-trends/internal/db/neo4j"
)

type NodeService interface {
	// FindNodeByName finds a node by its name
	FindNodeByName(ctx context.Context, name string) (*n4j.Node, error)

	// FindNodesByName finds all nodes with a given name
	FindNodesByName(ctx context.Context, name string) ([]n4j.Node, error)

	// FindTopTags finds the top n most popular tags based on node size
	FindTopTags(ctx context.Context, limit int) ([]n4j.Node, error)

	// FindTopTagsByDegreeCentrality finds the top n most popular tags based on degree centrality
	FindTopTagsByDegreeCentrality(ctx context.Context, limit int) ([]n4j.Node, error)

	// RecommendTags recommends tags based on the given tag and vagueness level
	RecommendTags(ctx context.Context, tag string, vagueness, limit int) ([]n4j.Node, error)

	// RecommendTagsSimple recommends tags based on the given tag
	RecommendTagsSimple(ctx context.Context, tag string, limit int) ([]string, error)
}

func NewSoTrendsService(repo n4j.NodeRepository) *SoTrendsService {
	return &SoTrendsService{Repository: repo}
}

func (s *SoTrendsService) FindNodeByName(ctx context.Context, name string) (*n4j.Node, error) {
	return s.Repository.FindNodeByName(ctx, name)
}

func (s *SoTrendsService) FindNodesByName(ctx context.Context, name string) ([]n4j.Node, error) {
	return s.Repository.FindNodesByName(ctx, name)
}

func (s *SoTrendsService) FindTopTags(ctx context.Context, limit int) ([]n4j.Node, error) {
	return s.Repository.FindTopTags(ctx, limit)
}

func (s *SoTrendsService) FindTopTagsByDegreeCentrality(ctx context.Context, limit int) ([]n4j.Node, error) {
	return s.Repository.FindTopTagsByDegreeCentrality(ctx, limit)
}

func (s *SoTrendsService) RecommendTags(ctx context.Context, tag string, vagueness, limit int) ([]n4j.Node, error) {
	return s.Repository.RecommendTags(ctx, tag, vagueness, limit)
}

func (s *SoTrendsService) RecommendTagsSimple(ctx context.Context, tag string, limit int) ([]string, error) {
	return s.Repository.RecommendTagsSimple(ctx, tag, limit)
}

type SoTrendsService struct {
	Repository n4j.NodeRepository
}
