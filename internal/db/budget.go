package db

import (
	"gorm.io/gorm"
	"bookkeeper-backend/internal/models"
)

type BudgetStore struct {
	DB *gorm.DB
}

func (s *BudgetStore) GetBudgetByID(id uint) (*models.Budget, error) {
	var b models.Budget
	if err := s.DB.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}
