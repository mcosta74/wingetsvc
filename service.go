package wingetsvc

import "context"

type Query struct {
	Term string `json:"term,omitempty"`
}

type Service interface {
	Search(ctx context.Context, query Query) ([]ServiceInfo, error)
	Versions(ctx context.Context, packageId string) ([]string, error)
}

func NewService(controller WingetController) Service {
	return &wingetSvc{
		controller: controller,
	}
}

type wingetSvc struct {
	controller WingetController
}

func (svc *wingetSvc) Search(ctx context.Context, query Query) ([]ServiceInfo, error) {
	return svc.controller.Search(ctx, query.Term)
}

func (svc *wingetSvc) Versions(ctx context.Context, packageId string) ([]string, error) {
	return svc.controller.Versions(ctx, packageId)
}
