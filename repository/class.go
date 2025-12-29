package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"changeme/types"
)

type ClassRepository struct {
	db *DB
}

func NewClassRepository(db *DB) *ClassRepository {
	return &ClassRepository{
		db: db,
	}
}

func (c *ClassRepository) createPeriod(period *types.Period) error {
	query := "INSERT INTO periods (name) VALUES (:name)"

	result, err := c.db.NamedExec(query, period)
	if err != nil {
		return fmt.Errorf("db named exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	period.ID = id

	return nil
}

func (c *ClassRepository) findPeriodByName(name string) (*types.Period, error) {
	var period types.Period

	query := "SELECT * FROM periods WHERE name = ?"
	if err := c.db.Get(&period, query, name); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &period, nil
}

func (c *ClassRepository) FindPeriodByID(id int64) (*types.Period, error) {
	var period types.Period

	query := "SELECT * FROM periods WHERE id = ?"
	if err := c.db.Get(&period, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &period, nil
}

func (c *ClassRepository) FindPeriods() ([]*types.Period, error) {
	var periods []*types.Period

	query := "SELECT * FROM periods"
	if err := c.db.Select(&periods, query); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	if len(periods) == 0 {
		return nil, sql.ErrNoRows
	}

	return periods, nil
}

func (c *ClassRepository) createQuota(quota *types.Quota) error {
	query := "INSERT INTO quotas (name) VALUES (:name)"

	result, err := c.db.NamedExec(query, quota)
	if err != nil {
		return fmt.Errorf("db named exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	quota.ID = id

	return nil
}

func (c *ClassRepository) findQuotaByName(name string) (*types.Quota, error) {
	var quota types.Quota

	query := "SELECT * FROM quotas WHERE name = ?"
	if err := c.db.Get(&quota, query, name); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	return &quota, nil
}

func (c *ClassRepository) findQuotaByID(id int64) (*types.Quota, error) {
	var quota types.Quota

	query := "SELECT * FROM quotas WHERE id = ?"
	if err := c.db.Get(&quota, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &quota, nil
}

func (c *ClassRepository) findOrCreatePeriod(period *types.Period) error {
	p, err := c.findPeriodByName(period.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find period: %w", err)
	} else if err == nil {
		period.ID = p.ID
		return nil
	}

	if err := c.createPeriod(period); err != nil {
		return fmt.Errorf("create period: %w", err)
	}

	return nil
}

func (c *ClassRepository) findOrCreateQuota(quota *types.Quota) error {
	q, err := c.findQuotaByName(quota.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find quota by name: %w", err)
	} else if err == nil {
		quota.ID = q.ID

		return nil
	}

	if err := c.createQuota(quota); err != nil {
		return fmt.Errorf("create quota: %w", err)
	}

	return nil
}

func (c *ClassRepository) findOrCreateClass(class *types.Class) error {
	cl, err := c.FindClassByPeriodAndQuota(class.Period.ID, class.Quota.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("find class by period and quota: %w", err)
	} else if err == nil {
		*class = *cl

		return nil
	}

	query := `
		INSERT INTO classes (
			period_id, 
			quota_id, 
			seats,
			minimum_score
		) VALUES (
			:period.id,
			:quota.id,
			:seats,
			:minimum_score
		)
	`

	result, err := c.db.NamedExec(query, class)
	if err != nil {
		return fmt.Errorf("db named exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}

	class.ID = id

	return nil
}

func (c *ClassRepository) CreateClass(class *types.Class) error {
	if err := c.findOrCreatePeriod(&class.Period); err != nil {
		return fmt.Errorf("find or create period: %w", err)
	}

	if err := c.findOrCreateQuota(&class.Quota); err != nil {
		return fmt.Errorf("find or create quota: %w", err)
	}

	if err := c.findOrCreateClass(class); err != nil {
		return err
	}

	return nil
}

func (c *ClassRepository) DeleteClass(id int64) error {
	query := "DELETE FROM classes WHERE ID = ?"

	_, err := c.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	return nil
}

func (c *ClassRepository) FindClassByPeriodAndQuota(periodID int64, quotaID int64) (*types.Class, error) {
	var class types.Class

	query := `
		SELECT classes.id, seats, minimum_score, periods.id as "period.id", periods.name as "period.name", quotas.id as "quota.id", quotas.name as "quota.name"
		FROM classes
		JOIN periods ON classes.period_id = periods.id
	    JOIN quotas ON classes.quota_id = quotas.id
		WHERE period_id = ? AND quota_id = ?
	`
	if err := c.db.Get(&class, query, periodID, quotaID); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &class, nil
}

func (c *ClassRepository) FindClassByID(id int64) (*types.Class, error) {
	var class types.Class

	query := `
		SELECT classes.id, seats, minimum_score, periods.id as "period.id", periods.name as "period.name", quotas.id as "quota.id", quotas.name as "quota.name"
		FROM classes
		JOIN periods ON classes.period_id = periods.id
	    JOIN quotas ON classes.quota_id = quotas.id
		WHERE classes.id = ?
	`
	if err := c.db.Get(&class, query, id); err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &class, nil
}

func (c *ClassRepository) FindClassesByPeriodID(periodID int64) ([]*types.Class, error) {
	var classes []*types.Class

	query := `
		SELECT classes.id, seats, minimum_score, periods.id as "period.id", periods.name as "period.name", quotas.id as "quota.id", quotas.name as "quota.name"
		FROM classes
		JOIN periods ON classes.period_id = periods.id
	    JOIN quotas ON classes.quota_id = quotas.id
		WHERE classes.period_id = ?
	`
	if err := c.db.Select(&classes, query, periodID); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	if len(classes) == 0 {
		return nil, sql.ErrNoRows
	}

	return classes, nil
}

func (c *ClassRepository) FindClasses() ([]*types.Class, error) {
	var classes []*types.Class

	query := `
		SELECT classes.id, seats, minimum_score, periods.id as "period.id", periods.name as "period.name", quotas.id as "quota.id", quotas.name as "quota.name"
		FROM classes
		JOIN periods ON classes.period_id = periods.id
	    JOIN quotas ON classes.quota_id = quotas.id
	`
	if err := c.db.Select(&classes, query); err != nil {
		return nil, fmt.Errorf("db select: %w", err)
	}

	if len(classes) == 0 {
		return nil, sql.ErrNoRows
	}

	return classes, nil
}
