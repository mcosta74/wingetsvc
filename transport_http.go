package wingetsvc

import (
	"context"
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

type contextKey int

var (
	loggerContextKey contextKey = 0
)

func MakeHTTPHandler(endpoints Endpoints, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(loggerToContext(logger), kithttp.PopulateRequestContext),
	}

	searchHandler := kithttp.NewServer(
		endpoints.Search,
		decodeSearchRequest,
		encodeResponse,
		// kithttp.EncodeJSONResponse,
		options...,
	)

	versionsHandler := kithttp.NewServer(
		endpoints.Versions,
		decodeVersionsRequest,
		encodeResponse,
		// kithttp.EncodeJSONResponse,
		options...,
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
		encodeError(ctx, failer.Failed(), rw)
		return nil
	}
	return kithttp.EncodeJSONResponse(ctx, rw, respose)
}

func encodeError(ctx context.Context, err error, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)

	logger := loggerFromContext(ctx)
	if logger != nil {
		logger.Error("error response", "err", err, "path", ctx.Value(kithttp.ContextKeyRequestPath))
	}

	tmp := map[string]string{
		"description": err.Error(),
	}
	_ = json.NewEncoder(rw).Encode(tmp)
}

func loggerToContext(logger *slog.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, loggerContextKey, logger)
	}
}

func loggerFromContext(ctx context.Context) *slog.Logger {
	val := ctx.Value(loggerContextKey)

	if logger, ok := val.(*slog.Logger); ok {
		return logger
	}
	return nil
}
