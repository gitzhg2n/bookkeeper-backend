package routes

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "bookkeeper-backend/middleware"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// AdminEntitlementHandler provides simple handlers to grant/revoke entitlements for testing.
type AdminEntitlementHandler struct {
    DB *gorm.DB
}

type entitlementRequest struct {
    UserID     uint   `json:"user_id"`
    FeatureKey string `json:"feature_key"`
    Enabled    bool   `json:"enabled"`
}

func NewAdminEntitlementHandler(db *gorm.DB) *AdminEntitlementHandler {
    return &AdminEntitlementHandler{DB: db}
}

// Upsert sets an entitlement on/off for a user. Protected to authenticated users with Role=="admin".
func (h *AdminEntitlementHandler) Upsert(w http.ResponseWriter, r *http.Request) {
    user, ok := middleware.UserFrom(r.Context())
    if !ok {
        writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
        return
    }
    if user.Role != "admin" {
        writeJSONError(r, w, "forbidden", http.StatusForbidden)
        return
    }

    var req entitlementRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSONError(r, w, "invalid payload", http.StatusBadRequest)
        return
    }

    enabledInt := 0
    if req.Enabled {
        enabledInt = 1
    }
    // Upsert using gorm clause for cross-db compatibility
    ent := map[string]interface{}{"user_id": req.UserID, "feature_key": req.FeatureKey, "enabled": enabledInt}
    if err := h.DB.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "feature_key"}}, DoUpdates: clause.AssignmentColumns([]string{"enabled", "updated_at"})}).Model(nil).Create(ent).Error; err != nil {
        writeJSONError(r, w, "db error", http.StatusInternalServerError)
        return
    }
    writeJSONSuccess(r, w, "ok", map[string]string{"status": "updated"})
}

// List entitlements for a user: /v1/admin/entitlements?user_id=123
func (h *AdminEntitlementHandler) List(w http.ResponseWriter, r *http.Request) {
    user, ok := middleware.UserFrom(r.Context())
    if !ok {
        writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
        return
    }
    if user.Role != "admin" {
        writeJSONError(r, w, "forbidden", http.StatusForbidden)
        return
    }
    q := r.URL.Query().Get("user_id")
    if q == "" {
        writeJSONError(r, w, "missing user_id", http.StatusBadRequest)
        return
    }
    id, err := strconv.Atoi(q)
    if err != nil {
        writeJSONError(r, w, "invalid user_id", http.StatusBadRequest)
        return
    }
    rows, err := h.DB.Raw(`SELECT feature_key, enabled, expires_at FROM entitlements WHERE user_id = ?`, id).Rows()
    if err != nil {
        writeJSONError(r, w, "db error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    var out []map[string]interface{}
    for rows.Next() {
        var fk string
        var enabled int
        var expires sql.NullString
        _ = rows.Scan(&fk, &enabled, &expires)
        m := map[string]interface{}{"feature_key": fk, "enabled": enabled == 1}
        if expires.Valid {
            m["expires_at"] = expires.String
        }
        out = append(out, m)
    }
    writeJSONSuccess(r, w, "ok", out)
}
