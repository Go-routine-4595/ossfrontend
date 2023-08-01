package service

import (
	"context"
	"github.com/Go-routine-4995/ossfrontend/domain"
)

type Broker interface {
	CreateRoutersRequest(r []domain.Router, tenant string) ([]byte, error)
	DeleteRoutersRequest(r []domain.Router, tenant string) error
	GetRoutersPage(paginationByte []byte, tenant string) (domain.Response, error)
	GetRouters(r domain.Router, tenant string) (domain.Router, error)
}

type IService interface {
}

type Service struct {
	broker Broker
}

func NewService(b interface{}) IService {
	return &Service{
		broker: b.(Broker),
	}
}

func (s *Service) CreateRouters(ctx context.Context, r []domain.Router, tenant string) ([]byte, error) {
	return s.broker.CreateRoutersRequest(r, tenant)
}

func (s *Service) DeleteRouters(ctx context.Context, r []domain.Router, tenant string) error {
	return s.broker.DeleteRoutersRequest(r, tenant)
}

func (s *Service) GetRoutersPage(ctx context.Context, paginationByte []byte, tenant string) (domain.Response, error) {
	return s.broker.GetRoutersPage(paginationByte, tenant)
}

func (s *Service) GetRouters(ctx context.Context, r domain.Router, tenant string) (domain.Router, error) {
	return s.broker.GetRouters(r, tenant)
}
