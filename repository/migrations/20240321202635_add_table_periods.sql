-- +goose Up
-- +goose StatementBegin
CREATE TABLE periods (
  id   INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_periods_name ON periods(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_periods_name;
DROP TABLE periods;
-- +goose StatementEnd
