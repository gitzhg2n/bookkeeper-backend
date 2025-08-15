package db

import (
	"gorm.io/gorm"
	"bookkeeper-backend/internal/models"
)

type UserSettingsStore struct {
	DB *gorm.DB
}

func (s *UserSettingsStore) GetByUserID(userID uint) (*models.UserSettings, error) {
	var us models.UserSettings
	if err := s.DB.Where("user_id = ?", userID).First(&us).Error; err != nil {
		return nil, err
	}
	return &us, nil
}

func (s *UserSettingsStore) Upsert(userID uint, threshold int64) error {
	var us models.UserSettings
	if err := s.DB.Where("user_id = ?", userID).First(&us).Error; err == nil {
		us.LargeTransactionThreshold = threshold
		return s.DB.Save(&us).Error
	}
	us = models.UserSettings{UserID: userID, LargeTransactionThreshold: threshold}
	return s.DB.Create(&us).Error
}
