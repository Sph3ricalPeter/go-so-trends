package service

import (
	"context"

	"github.com/Sph3ricalPeter/go-so-trends/internal/db/neo4j"
)

type SoTrendsService struct {
	Repository neo4j.SoNodeRepository
}

// find top 5 most popular tags based on node size
func (s *SoTrendsService) FindTopTags(ctx context.Context) ([]neo4j.Node, error) {
	return s.Repository.FindTopTags(ctx)
}
