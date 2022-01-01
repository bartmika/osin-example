package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bartmika/osin-example/internal/models"
)

type TenantRepo struct {
	db *sql.DB
}

func NewTenantRepo(db *sql.DB) *TenantRepo {
	return &TenantRepo{
		db: db,
	}
}

func (r *TenantRepo) Insert(ctx context.Context, m *models.Tenant) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO tenants (
        uuid, name, state, timezone, language, created_time, modified_time
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    )
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println("Insert")
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.UUID, m.Name, m.State, m.Timezone, m.Language, m.CreatedTime, m.ModifiedTime,
	)
	return err
}

func (r *TenantRepo) UpdateByID(ctx context.Context, m *models.Tenant) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        tenants
    SET
        name = $1,
		state = $2,
		timezone = $3,
		language = $4,
		created_time = $5,
		modified_time = $6
    WHERE
        id = $7
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, m.Name, m.State, m.Timezone, m.Language, m.CreatedTime, m.ModifiedTime, m.ID,
	)
	return err
}

func (r *TenantRepo) GetByID(ctx context.Context, id uint64) (*models.Tenant, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Tenant)

	query := `
    SELECT
        id, uuid, name, state, timezone, language, created_time, modified_time
    FROM
        tenants
    WHERE
        id = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.UUID, &m.Name, &m.State, &m.Timezone, &m.Language, &m.CreatedTime, &m.ModifiedTime,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *TenantRepo) GetByName(ctx context.Context, name string) (*models.Tenant, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Tenant)

	query := `
    SELECT
        id, uuid, name, state, timezone, language, created_time, modified_time
    FROM
        tenants
    WHERE
        name = $1
    `
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&m.ID, &m.UUID, &m.Name, &m.State, &m.Timezone, &m.Language, &m.CreatedTime, &m.ModifiedTime,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *TenantRepo) CheckIfExistsByID(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        tenants
    WHERE
        id = $1
    `

	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *TenantRepo) CheckIfExistsByName(ctx context.Context, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        tenants
    WHERE
        name = $1
    `

	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *TenantRepo) InsertOrUpdateByID(ctx context.Context, m *models.Tenant) error {
	if m.ID == 0 {
		return r.Insert(ctx, m)
	}

	doesExist, err := r.CheckIfExistsByID(ctx, m.ID)
	if err != nil {
		return err
	}

	if doesExist == false {
		return r.Insert(ctx, m)
	}
	return r.UpdateByID(ctx, m)
}
