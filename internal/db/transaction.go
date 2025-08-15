package db

import (
	"gorm.io/gorm"
	"bookkeeper-backend/internal/models"
)

type TransactionStore struct {
	DB *gorm.DB
}

func (s *TransactionStore) GetTransactionByID(id uint) (*models.Transaction, error) {
	var t models.Transaction
	if err := s.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
