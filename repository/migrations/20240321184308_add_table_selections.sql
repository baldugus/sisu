-- +goose Up
-- +goose StatementBegin
CREATE TABLE selections (
  id          INTEGER PRIMARY KEY,
  kind        INTEGER UNIQUE NOT NULL,
  name        TEXT NOT NULL,
  date        TEXT NOT NULL,
  institution TEXT NOT NULL,
  course      TEXT NOT NULL
);

CREATE INDEX idx_selections_kind ON selections(kind);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_selections_kind;
DROP TABLE selections;
-- +goose StatementEnd
