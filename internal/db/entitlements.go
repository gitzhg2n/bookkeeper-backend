package db

import (
    "fmt"

    "gorm.io/gorm"
)

// UserHasEntitlement checks whether the user has an enabled entitlement for a featureKey.
// Uses a simple SQL query compatible with sqlite and postgres via gorm.DB.
func UserHasEntitlement(gdb *gorm.DB, userID uint, featureKey string) (bool, error) {
    var count int64
    if err := gdb.Raw(`SELECT COUNT(1) FROM entitlements WHERE user_id = ? AND feature_key = ? AND enabled = 1 AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)`, userID, featureKey).Scan(&count).Error; err != nil {
        return false, fmt.Errorf("query entitlements: %w", err)
    }
    return count > 0, nil
}
