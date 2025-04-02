package service

import (
	"context"
	"lesson3/internal/models"
	test "lesson3/pkg/api/test/api"

	"github.com/google/uuid"
)

type Storage interface {
	CreateOrder(order models.Order) (string, error)
	GetOrder(id string) (models.Order, error)
	UpdateOrder(order models.Order) (models.Order, error)
	DeleteOrder(id string) error
	ListOrders() map[string]models.Order
}

type Service struct {
	test.OrderServiceServer
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateOrder(ctx context.Context, req *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	order := models.Order{
		ID:       uuid.New().String(),
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	id, err := s.storage.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	return &test.CreateOrderResponse{
		Id: id,
	}, nil
}

func (s *Service) GetOrder(ctx context.Context, req *test.GetOrderRequest) (*test.GetOrderResponse, error) {
	order, err := s.storage.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return &test.GetOrderResponse{
		Order: &test.Order{
			Id:       order.ID,
			Item:     order.Item,
			Quantity: order.Quantity,
		},
	}, nil
}

func (s *Service) ListOrders(ctx context.Context, req *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {
	orders := s.storage.ListOrders()

	var orders_test []*test.Order

	for _, order := range orders {
		orders_test = append(orders_test, &test.Order{
			Id:       order.ID,
			Item:     order.Item,
			Quantity: order.Quantity,
		})
	}

	return &test.ListOrdersResponse{
		Orders: orders_test,
	}, nil
}

func (s *Service) UpdateOrder(ctx context.Context, req *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	order := models.Order{
		ID:       req.Id,
		Item:     req.Item,
		Quantity: req.Quantity,
	}

	order, err := s.storage.UpdateOrder(order)
	if err != nil {
		return nil, err
	}

	order_test := test.Order{
		Id:       order.ID,
		Item:     order.Item,
		Quantity: order.Quantity,
	}

	return &test.UpdateOrderResponse{
		Order: &order_test,
	}, nil
}

func (s *Service) DeleteOrder(ctx context.Context, req *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {
	err := s.storage.DeleteOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return &test.DeleteOrderResponse{
		Success: true,
	}, nil
}
