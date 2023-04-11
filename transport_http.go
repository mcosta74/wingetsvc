package wingetsvc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(endpoints Endpoints) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// options := []kithttp.ServerOption{

	// }

	searchHandler := kithttp.NewServer(
		endpoints.Search,
		decodeSearchRequest,
		encodeResponse,
		// kithttp.EncodeJSONResponse,
	)

	versionsHandler := kithttp.NewServer(
		endpoints.Versions,
		decodeVersionsRequest,
		encodeResponse,
		// kithttp.EncodeJSONResponse,
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

func encodeResponse(ctx context.Context, rw http.ResponseWriter, respose any) error {
	if failer, ok := respose.(endpoint.Failer); ok && failer.Failed() != nil {
		return encodeError(ctx, rw, failer.Failed())
	}
	return kithttp.EncodeJSONResponse(ctx, rw, respose)
}

func encodeError(ctx context.Context, rw http.ResponseWriter, err error) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)

	tmp := map[string]string{
		"description": err.Error(),
	}
	return json.NewEncoder(rw).Encode(tmp)
}
