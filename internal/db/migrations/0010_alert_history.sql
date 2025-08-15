-- +migrate Up
CREATE TABLE IF NOT EXISTS alert_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alert_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    triggered_at DATETIME NOT NULL,
    details TEXT,
    FOREIGN KEY(alert_id) REFERENCES investment_alerts(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +migrate Down
DROP TABLE IF EXISTS alert_history;
