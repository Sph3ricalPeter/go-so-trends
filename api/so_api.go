package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sph3ricalPeter/go-so-trends/internal/service"
)

type NodeApi interface {
	// MountRoutes mounts the API routes
	MountRoutes() http.Handler

	// Alive checks if the API is alive
	Alive(w http.ResponseWriter, r *http.Request)

	// FindNodeByName finds a node by its name
	FindNodeByName(w http.ResponseWriter, r *http.Request)

	// FindTopTags finds the top 5 most popular tags based on node size
	FindTopTags(w http.ResponseWriter, r *http.Request)

	// FindTopTagsByDegreeCentrality finds the top 5 most popular tags based on degree centrality
	FindTopTagsByDegreeCentrality(w http.ResponseWriter, r *http.Request)

	// RecommendTags recommends tags based on the given tag
	RecommendTags(w http.ResponseWriter, r *http.Request)

	// RecommendTagsSimple recommends tags based on the given tag
	RecommendTagsSimple(w http.ResponseWriter, r *http.Request)
}

func NewSoTrendsApi(ctx context.Context, service service.NodeService) *SoTrendsApi {
	return &SoTrendsApi{Service: service, Context: context.Background()}
}

func (api *SoTrendsApi) MountRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/alive", api.Alive)
	mux.HandleFunc("/find/{name}", api.FindNodeByName)
	mux.HandleFunc("/toptags/size", api.FindTopTags)
	mux.HandleFunc("/toptags/dc", api.FindTopTagsByDegreeCentrality)
	mux.HandleFunc("/recommend/{tag}/{vagueness}", api.RecommendTags)
	mux.HandleFunc("/recommend/{tag}", api.RecommendTagsSimple)

	return ContentTypeMiddleware(mux)
}

func (api *SoTrendsApi) Alive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I'm alive!"))
}

func (api *SoTrendsApi) FindNodeByName(w http.ResponseWriter, r *http.Request) {
	data, err := api.Service.FindNodeByName(api.Context, r.PathValue("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (api *SoTrendsApi) FindTopTags(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	data, err := api.Service.FindTopTags(api.Context, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (api *SoTrendsApi) FindTopTagsByDegreeCentrality(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	data, err := api.Service.FindTopTagsByDegreeCentrality(api.Context, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (api *SoTrendsApi) RecommendTags(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	vagueness, err := strconv.Atoi(r.PathValue("vagueness"))
	if err != nil {
		http.Error(w, "Invalid vagueness", http.StatusBadRequest)
		return
	}
	data, err := api.Service.RecommendTags(api.Context, tag, vagueness, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (api *SoTrendsApi) RecommendTagsSimple(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	data, err := api.Service.RecommendTagsSimple(api.Context, tag, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// ContentTypeMiddleware sets the Content-Type header to application/json
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type SoTrendsApi struct {
	Service service.NodeService
	Context context.Context
}
