-- +goose Up
-- +goose StatementBegin
CREATE TABLE classes (
  id            INTEGER PRIMARY KEY,
  period_id     INTEGER NOT NULL REFERENCES periods ON DELETE RESTRICT ON UPDATE RESTRICT,
  quota_id      INTEGER NOT NULL REFERENCES quotas ON DELETE RESTRICT ON UPDATE RESTRICT,
  seats         INTEGER NOT NULL,
  minimum_score REAL NOT NULL
);

CREATE UNIQUE INDEX idx_classes_period_quota ON classes(period_id, quota_id);
CREATE INDEX idx_classes_period ON classes(period_id);
CREATE INDEX idx_classes_quota ON classes(quota_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_classes_period_quota;
DROP INDEX idx_classes_period;
DROP INDEX idx_classes_quota;
DROP TABLE classes;
-- +goose StatementEnd
