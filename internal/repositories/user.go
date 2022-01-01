package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/bartmika/osin-example/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Insert(ctx context.Context, m *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO users (
        uuid, tenant_id, email, first_name, last_name, password_algorithm,
		password_hash, state, role_id, timezone, language, created_time, modified_time,
		joined_time, salt, was_email_activated, pr_access_code, pr_expiry_time,
		name, lexical_name
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
    )`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.UUID, m.TenantID, m.Email, m.FirstName, m.LastName, m.PasswordAlgorithm,
		m.PasswordHash, m.State, m.RoleID, m.Timezone, m.Language, m.CreatedTime, m.ModifiedTime,
		m.JoinedTime, m.Salt, m.WasEmailActivated, m.PrAccessCode, m.PrExpiryTime,
		m.Name, m.LexicalName,
	)
	return err
}

func (r *UserRepo) UpdateByID(ctx context.Context, m *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        users
    SET
        tenant_id = $1,
		email = $2,
		first_name = $3,
		last_name = $4,
		password_algorithm = $5,
		password_hash = $6,
		state = $7,
		role_id = $8,
		timezone = $9,
		language = $10
		created_time = $11,
		modified_time = $12,
		joined_time = $13,
		salt = $14,
		was_email_activated = $15,
		pr_access_code = $16,
		pr_expiry_time = $17,
		name = $18,
		lexical_name = $19
    WHERE
        id = $20`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantID,
		m.Email,
		m.FirstName,
		m.LastName,
		m.PasswordAlgorithm,
		m.PasswordHash,
		m.State,
		m.RoleID,
		m.Timezone,
		m.Language,
		m.CreatedTime,
		m.ModifiedTime,
		m.JoinedTime,
		m.Salt,
		m.WasEmailActivated,
		m.PrAccessCode,
		m.PrExpiryTime,
		m.Name,
		m.LexicalName,
		m.ID,
	)
	return err
}

func (r *UserRepo) UpdateByEmail(ctx context.Context, m *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE
        users
    SET
        tenant_id = $1,
		first_name = $2,
		last_name = $3,
		password_algorithm = $4,
		password_hash = $5,
		state = $6,
		role_id = $7,
		timezone = $8,
		language = $9,
		created_time = $10,
		modified_time = $11,
		joined_time = $12,
		salt = $13,
		was_email_activated = $14,
		pr_access_code = $15,
		pr_expiry_time = $16,
		name = $17,
		lexical_name = $18
    WHERE
        email = $19`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		m.TenantID,
		m.FirstName,
		m.LastName,
		m.PasswordAlgorithm,
		m.PasswordHash,
		m.State,
		m.RoleID,
		m.Timezone,
		m.Language,
		m.CreatedTime,
		m.ModifiedTime,
		m.JoinedTime,
		m.Salt,
		m.WasEmailActivated,
		m.PrAccessCode,
		m.PrExpiryTime,
		m.Name,
		m.LexicalName,
		m.Email,
	)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.User)

	query := `
    SELECT
		uuid, id, tenant_id, email, first_name, last_name, password_algorithm,
		password_hash, state, role_id, timezone, language, created_time, modified_time,
		joined_time, salt, was_email_activated, pr_access_code, pr_expiry_time,
		name, lexical_name
    FROM
        users
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.UUID, &m.ID, &m.TenantID, &m.Email, &m.FirstName, &m.LastName, &m.PasswordAlgorithm,
		&m.PasswordHash, &m.State, &m.RoleID, &m.Timezone, &m.Language, &m.CreatedTime, &m.ModifiedTime,
		&m.JoinedTime, &m.Salt, &m.WasEmailActivated, &m.PrAccessCode, &m.PrExpiryTime,
		&m.Name, &m.LexicalName,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.User)

	query := `
    SELECT
		uuid, id, tenant_id, email, first_name, last_name, password_algorithm,
		password_hash, state, role_id, timezone, language, created_time, modified_time,
		joined_time, salt, was_email_activated, pr_access_code, pr_expiry_time,
		name, lexical_name
    FROM
        users
    WHERE
        email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&m.UUID, &m.ID, &m.TenantID, &m.Email, &m.FirstName, &m.LastName, &m.PasswordAlgorithm,
		&m.PasswordHash, &m.State, &m.RoleID, &m.Timezone, &m.Language, &m.CreatedTime, &m.ModifiedTime,
		&m.JoinedTime, &m.Salt, &m.WasEmailActivated, &m.PrAccessCode, &m.PrExpiryTime,
		&m.Name, &m.LexicalName,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *UserRepo) GetByUUID(ctx context.Context, uid string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m := new(models.User)

	query := `
    SELECT
        uuid, id, tenant_id, email, first_name, last_name, password_algorithm,
		password_hash, state, role_id, timezone, language, created_time, modified_time,
		joined_time, salt, was_email_activated, pr_access_code, pr_expiry_time,
		name, lexical_name
    FROM
        users
    WHERE
        uuid = $1`
	err := r.db.QueryRowContext(ctx, query, uid).Scan(
		&m.UUID, &m.ID, &m.TenantID, &m.Email, &m.FirstName, &m.LastName, &m.PasswordAlgorithm,
		&m.PasswordHash, &m.State, &m.RoleID, &m.Timezone, &m.Language, &m.CreatedTime, &m.ModifiedTime,
		&m.JoinedTime, &m.Salt, &m.WasEmailActivated, &m.PrAccessCode, &m.PrExpiryTime,
		&m.Name, &m.LexicalName,
	)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return nil, nil
		} else { // CASE 2 OF 2: All other errors.
			return nil, err
		}
	}
	return m, nil
}

func (r *UserRepo) CheckIfExistsByID(ctx context.Context, id uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        users
    WHERE
        id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *UserRepo) CheckIfExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool

	query := `
    SELECT
        1
    FROM
        users
    WHERE
        email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		// CASE 1 OF 2: Cannot find record with that email.
		if err == sql.ErrNoRows {
			return false, nil
		} else { // CASE 2 OF 2: All other errors.
			return false, err
		}
	}
	return exists, nil
}

func (r *UserRepo) InsertOrUpdateByID(ctx context.Context, m *models.User) error {
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

func (r *UserRepo) InsertOrUpdateByEmail(ctx context.Context, m *models.User) error {
	if m.ID == 0 {
		return r.Insert(ctx, m)
	}

	doesExist, err := r.CheckIfExistsByEmail(ctx, m.Email)
	if err != nil {
		return err
	}

	if doesExist == false {
		return r.Insert(ctx, m)
	}
	return r.UpdateByEmail(ctx, m)
}
