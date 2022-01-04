package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bartmika/osin-example/internal/models"
)

type ApplicationRepo struct {
	db *sql.DB
}

func NewApplicationRepo(db *sql.DB) *ApplicationRepo {
	return &ApplicationRepo{
		db: db,
	}
}

func (r *ApplicationRepo) Insert(ctx context.Context, m *models.Application) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO applications (
        uuid, tenant_id, name, description, scope, redirect_url, image_url, state, client_id, client_secret
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
    )
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println("ApplicationRepo|Insert|err", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.UUID, m.TenantID, m.Name, m.Description, m.Scope, m.RedirectURL, m.ImageURL, m.State, m.ClientID, m.ClientSecret,
	)
	return err
}

func (r *ApplicationRepo) UpdateByID(ctx context.Context, m *models.Application) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        applications
    SET
        name = $1,
		description = $2,
		scope = $3,
		redirect_url = $4,
		image_url = $5,
		state = $6,
		client_id = $7,
		client_secret = $8,
		created_time = $9,
		modified_time = $10
    WHERE
        id = $11
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, m.Name, m.Description, m.Scope, m.RedirectURL, m.ImageURL, m.State,
		m.ClientID, m.ClientSecret, m.CreatedTime, m.ModifiedTime, m.ID,
	)
	return err
}

func (r *ApplicationRepo) GetByID(ctx context.Context, id uint64) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Application)

	query := `
    SELECT
        id, uuid, tenant_id, name, description, scope, redirect_url, image_url, state,
		client_id, client_secret, created_time, modified_time
    FROM
        applications
    WHERE
        id = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.UUID, &m.TenantID, &m.Name, &m.Description, &m.Scope,
		&m.RedirectURL, &m.ImageURL, &m.State, &m.ClientID, &m.ClientSecret,
		&m.CreatedTime, &m.ModifiedTime,
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

func (r *ApplicationRepo) GetByUUID(ctx context.Context, id string) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Application)

	query := `
    SELECT
        id, uuid, tenant_id, name, description, scope, redirect_url, image_url,
		state, client_id, client_secret, created_time, modified_time
    FROM
        applications
    WHERE
        uuid = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.UUID, &m.TenantID, &m.Name, &m.Description, &m.Scope,
		&m.RedirectURL, &m.ImageURL, &m.State, &m.ClientID, &m.ClientSecret,
		&m.CreatedTime, &m.ModifiedTime,
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

func (r *ApplicationRepo) GetByClientID(ctx context.Context, cid string) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.Application)

	query := `
    SELECT
        id, uuid, tenant_id, name, description, scope, redirect_url, image_url,
		state, client_id, client_secret, created_time, modified_time
    FROM
        applications
    WHERE
        client_id = $1
    `
	err := r.db.QueryRowContext(ctx, query, cid).Scan(
		&m.ID, &m.UUID, &m.TenantID, &m.Name, &m.Description, &m.Scope,
		&m.RedirectURL, &m.ImageURL, &m.State, &m.ClientID, &m.ClientSecret,
		&m.CreatedTime, &m.ModifiedTime,
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

func (r *ApplicationRepo) CheckIfExistsByID(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        applications
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

func (r *ApplicationRepo) CheckIfRunningByClientID(ctx context.Context, clientID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        applications
    WHERE
        client_id = $1
	AND
	    state = $2
    `

	err := r.db.QueryRowContext(ctx, query, clientID, models.ApplicationRunningState).Scan(&exists)
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

func (r *ApplicationRepo) InsertOrUpdateByID(ctx context.Context, m *models.Application) error {
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

func (s *ApplicationRepo) DeleteByID(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM applications WHERE id = $1;`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows { // CASE 1 OF 2: Cannot find record with that ID.
			return nil
		}
		// CASE 2 OF 2: All other errors.
		return err
	}
	return nil
}
