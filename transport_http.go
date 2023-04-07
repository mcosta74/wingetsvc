package wingetsvc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(endpoints Endpoints) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	searchHandler := kithttp.NewServer(
		endpoints.Search,
		decodeSearchRequest,
		kithttp.EncodeJSONResponse,
	)

	versionsHandler := kithttp.NewServer(
		endpoints.Versions,
		decodeVersionsRequest,
		kithttp.EncodeJSONResponse,
	)

	r.Route("/api", func(r chi.Router) {
		r.Method(http.MethodPost, "/search", searchHandler)
		r.Method(http.MethodPost, "/versions", versionsHandler)
	})

	return r
}

func decodeSearchRequest(ctx context.Context, r *http.Request) (any, error) {
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeVersionsRequest(ctx context.Context, r *http.Request) (any, error) {
	var request versionsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
