-- +goose Up
-- +goose StatementBegin
CREATE TABLE rollcalls (
  id     INTEGER PRIMARY KEY,
  number INTEGER UNIQUE NOT NULL,
  status TEXT CHECK ( status IN ('CALLING', 'DONE') )
);

CREATE INDEX idx_rollcalls_number ON rollcalls(number);
CREATE INDEX idx_rollcalls_status ON rollcalls(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_rollcalls_number;
DROP INDEX idx_rollcalls_status;
DROP TABLE rollcalls;
-- +goose StatementEnd
