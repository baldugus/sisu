-- +goose Up
-- +goose StatementBegin
CREATE TABLE applications (
  id                     INTEGER PRIMARY KEY,
  enrollment_id          TEXT NOT NULL UNIQUE,
  class_id               INTEGER NOT NULL REFERENCES classes ON DELETE RESTRICT ON UPDATE RESTRICT,
  option                 INTEGER NOT NULL,
  languages_score        REAL NOT NULL,
  humanities_score       REAL NOT NULL,
  natural_sciences_score REAL NOT NULL,
  mathematics_score      REAL NOT NULL,
  essay_score            REAL NOT NULL,
  composite_score        REAL NOT NULL,
  ranking                INTEGER NOT NULL,
  status                 TEXT CHECK (status IN ('APPROVED', 'WAITING', 'ABSENT', 'ENROLLED')),
  selection_id           INTEGER NOT NULL REFERENCES selections ON DELETE RESTRICT ON UPDATE RESTRICT,
  rollcall_id            INTEGER DEFAULT NULL REFERENCES rollcalls ON DELETE SET NULL ON UPDATE RESTRICT,
  applicant_id           INTEGER NOT NULL REFERENCES applicants ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE INDEX idx_applications_selection_id ON applications(selection_id);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_rollcall_id ON applications(rollcall_id);
CREATE INDEX idx_applications_class_id ON applications(class_id);
CREATE INDEX idx_applications_applicant_id ON applications(applicant_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_applications_selection_id;
DROP INDEX idx_applications_status;
DROP INDEX idx_applications_rollcall_id;
DROP INDEX idx_applications_class_id;
DROP INDEX idx_applications_applicant_id;
DROP TABLE applications;
-- +goose StatementEnd
