package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bartmika/osin-example/internal/models"
)

type AuthorizedApplicationRepo struct {
	db *sql.DB
}

func NewAuthorizedApplicationRepo(db *sql.DB) *AuthorizedApplicationRepo {
	return &AuthorizedApplicationRepo{
		db: db,
	}
}

func (r *AuthorizedApplicationRepo) Insert(ctx context.Context, m *models.AuthorizedApplication) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO authorized_applications (
        uuid, tenant_id, application_id, user_id, state,
		created_time, modified_time
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
		m.UUID, m.TenantID, m.ApplicationID, m.UserID, m.State,
		m.CreatedTime, m.ModifiedTime,
	)
	return err
}

func (r *AuthorizedApplicationRepo) UpdateByID(ctx context.Context, m *models.AuthorizedApplication) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        authorized_applications
    SET
        application_id = $1,
		user_id = $2,
		state = $3,
		created_time = $4,
		modified_time = $5
    WHERE
        id = $6
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, m.ApplicationID, m.UserID, m.State,
		m.CreatedTime, m.ModifiedTime, m.ID,
	)
	return err
}

func (r *AuthorizedApplicationRepo) GetByID(ctx context.Context, id uint64) (*models.AuthorizedApplication, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AuthorizedApplication)

	query := `
    SELECT
        id, uuid, tenant_id, application_id, user_id, state,
		created_time, modified_time
    FROM
        authorized_applications
    WHERE
        id = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.UUID, &m.TenantID, &m.ApplicationID, &m.UserID,
		&m.State, &m.CreatedTime, &m.ModifiedTime,
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

func (r *AuthorizedApplicationRepo) GetByUUID(ctx context.Context, id string) (*models.AuthorizedApplication, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AuthorizedApplication)

	query := `
    SELECT
        id, uuid, tenant_id, application_id, user_id, state,
		created_time, modified_time
    FROM
        authorized_applications
    WHERE
        uuid = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.UUID, &m.TenantID, &m.ApplicationID, &m.UserID,
		&m.State, &m.CreatedTime, &m.ModifiedTime,
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

func (r *AuthorizedApplicationRepo) GetByUserIDAndApplicationID(ctx context.Context, uid uint64, aid uint64) (*models.AuthorizedApplication, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.AuthorizedApplication)

	query := `
    SELECT
        id, uuid, tenant_id, application_id, user_id, state,
		created_time, modified_time
    FROM
        authorized_applications
    WHERE
        user_id = $1
	AND
	    application_id = $2
    `
	err := r.db.QueryRowContext(ctx, query, uid, aid).Scan(
		&m.ID, &m.UUID, &m.TenantID, &m.ApplicationID, &m.UserID,
		&m.State, &m.CreatedTime, &m.ModifiedTime,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		// CASE 2 OF 2: All other errors.
		return nil, err
	}
	return m, nil
}

func (r *AuthorizedApplicationRepo) CheckIfExistsByID(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        authorized_applications
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

func (r *AuthorizedApplicationRepo) CheckIfPermissionGrantedByUserIDAndByApplicationID(ctx context.Context, uid uint64, aid uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        authorized_applications
    WHERE
        user_id = $1
	AND
		application_id = $2
	AND
		state = $3
    `

	err := r.db.QueryRowContext(ctx, query, uid, aid, models.AuthorizedApplicationPermissionGrantedState).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that id.
		if err == sql.ErrNoRows {
			return false, nil
		}
		// CASE 2 OF 2: All other errors.
		return false, err

	}
	return exists, nil
}

func (r *AuthorizedApplicationRepo) InsertOrUpdateByID(ctx context.Context, m *models.AuthorizedApplication) error {
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

func (s *AuthorizedApplicationRepo) DeleteByID(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM authorized_applications WHERE id = $1;`
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
