-- +goose Up
-- +goose StatementBegin
CREATE TABLE quotas (
  id   INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_quotas_name ON quotas(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_quotas_name;
DROP TABLE quotas;
-- +goose StatementEnd
