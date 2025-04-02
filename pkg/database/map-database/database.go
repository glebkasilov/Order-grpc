package database

import (
	"errors"
	models "lesson3/internal/models"
	"sync"
)

type Database struct {
	mu      sync.Mutex
	Storage map[string]models.Order
}

func New() *Database {
	return &Database{
		Storage: make(map[string]models.Order),
	}
}

func (s *Database) CreateOrder(order models.Order) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Storage[order.ID] = order
	return order.ID, nil
}

func (s *Database) GetOrder(id string) (models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if (s.Storage[id] == models.Order{}) {
		return models.Order{}, errors.New("order not found")
	}

	return s.Storage[id], nil
}

func (s *Database) ListOrders() map[string]models.Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.Storage
}

func (s *Database) UpdateOrder(order models.Order) (models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if (s.Storage[order.ID] == models.Order{}) {
		return models.Order{}, errors.New("order not found")
	}

	s.Storage[order.ID] = order
	return order, nil
}

func (s *Database) DeleteOrder(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if (s.Storage[id] == models.Order{}) {
		return errors.New("order not found")
	}

	delete(s.Storage, id)

	return nil
}
