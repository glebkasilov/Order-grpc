package postgres

import (
	"fmt"
	"lesson3/internal/config"
	"lesson3/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New() (*Storage, error) {
	dbconfig, err := config.LoadDatabase()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbconfig.Host, dbconfig.Username, dbconfig.Password, dbconfig.Database, dbconfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Order{})

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (s *Storage) CreateOrder(order models.Order) (string, error) {
	if err := s.db.Create(&order).Error; err != nil {
		return "", err
	}

	return order.ID, nil
}

func (s *Storage) GetOrder(id string) (models.Order, error) {
	var order models.Order
	if err := s.db.Where("id = ?", id).First(&order).Error; err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (s *Storage) UpdateOrder(order models.Order) (models.Order, error) {
	if err := s.db.Save(&order).Error; err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (s *Storage) ListOrders() map[string]models.Order {
	var orders []models.Order
	if err := s.db.Find(&orders).Error; err != nil {
		return nil
	}

	orderMap := make(map[string]models.Order)
	for _, order := range orders {
		orderMap[order.ID] = order
	}

	return orderMap
}

func (s *Storage) DeleteOrder(id string) error {
	if err := s.db.Delete(&models.Order{}, id).Error; err != nil {
		return err
	}

	return nil
}
