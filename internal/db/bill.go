package db

import (
	"time"
	"gorm.io/gorm"
	"bookkeeper-backend/internal/models"
)

type BillStore struct {
	DB *gorm.DB
}

func (s *BillStore) ListDueInDays(userID uint, days int) ([]models.Bill, error) {
	var bills []models.Bill
	now := time.Now()
	future := now.Add(time.Duration(days) * 24 * time.Hour)
	if err := s.DB.Where("user_id = ? AND next_due >= ? AND next_due <= ?", userID, now, future).Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}
