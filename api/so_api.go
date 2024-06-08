package api

import (
	"context"
	"encoding/json"
	"net/http"

	n4j "github.com/Sph3ricalPeter/go-so-trends/internal/db/neo4j"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type SoTrendsApi struct {
	Repository n4j.NodeRepository
	Context    context.Context
}

func (api *SoTrendsApi) MountRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/alive", api.Alive)
	mux.HandleFunc("/find/{name}", api.FindNodeByName)
	mux.HandleFunc("/toptags", api.FindTopTags)

	return ContentTypeMiddleware(mux)
}

func (api *SoTrendsApi) Alive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I'm alive!"))
}

func (api *SoTrendsApi) FindNodeByName(w http.ResponseWriter, r *http.Request) {
	data, err := api.Repository.FindNodeByName(api.Context, r.PathValue("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (api *SoTrendsApi) FindTopTags(w http.ResponseWriter, r *http.Request) {
	data, err := api.Repository.FindTopTags(api.Context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
