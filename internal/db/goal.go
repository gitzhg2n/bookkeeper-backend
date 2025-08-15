package db

import (
	"gorm.io/gorm"
	"bookkeeper-backend/internal/models"
)

type GoalStore struct {
	DB *gorm.DB
}

func (s *GoalStore) ListByUser(userID uint) ([]models.Goal, error) {
	var goals []models.Goal
	if err := s.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}
