package services

import (
	"context"
	"errors"

	"microservice-test-order-go/models"
	"microservice-test-order-go/repositories"
)

type OrderService struct {
	repo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Validación básica antes de guardar
	if order.Name == "" || order.Email == "" {
		return errors.New("name and email are required")
	}

	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.repo.FindAll(ctx)
}

func (s *OrderService) GetOrderById(ctx context.Context, id string) (*models.Order, error) {
	order, err := s.repo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, id string, order *models.Order) error {
	// Verificar que los campos principales no sean vacíos
	if order.Name == "" || order.Email == "" {
		return errors.New("name and email are required")
	}

	return s.repo.Update(ctx, id, order)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
