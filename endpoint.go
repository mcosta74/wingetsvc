package wingetsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Search   endpoint.Endpoint
	Versions endpoint.Endpoint
}

type searchRequest struct {
	Query Query `json:"query"`
}

type searchResponse struct {
	V   []ServiceInfo `json:"v"`
	Err string        `json:"err,omitempty"`
}

type versionsRequest struct {
	PackageID string `json:"package_id"`
}

type versionsResponse struct {
	V   []string `json:"v"`
	Err string   `json:"err,omitempty"`
}

func MakeEndpoints(svc Service) Endpoints {
	return Endpoints{
		Search:   makeSearchEndpoint(svc),
		Versions: makeVersionsEndpoint(svc),
	}
}

func makeSearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req := request.(searchRequest)

		v, err := svc.Search(ctx, req.Query)
		if err != nil {
			return searchResponse{v, err.Error()}, nil
		}
		return searchResponse{v, ""}, nil
	}
}

func makeVersionsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req := request.(versionsRequest)

		v, err := svc.Versions(ctx, req.PackageID)
		if err != nil {
			return versionsResponse{v, err.Error()}, nil
		}
		return versionsResponse{v, ""}, nil
	}
}
