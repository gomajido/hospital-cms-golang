package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/google/uuid"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password, name, phone, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	user.ID = uuid.New().String()
	user.Status = "inactive" // Default status

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Phone,
		user.Status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, name, phone, status, email_verified_at, created_at, updated_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Status,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Get user roles
	roles, err := r.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return user, nil
}

func (r *authRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password, name, phone, status, email_verified_at, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	user := &domain.User{}
	var deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Status,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	// Get user roles
	roles, err := r.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return user, nil
}

func (r *authRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = ?, phone = ?, status = ?, email_verified_at = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Phone,
		user.Status,
		user.EmailVerifiedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *authRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *authRepository) GetRolesByNames(ctx context.Context, names []string) ([]domain.Role, error) {
	// Create placeholders for the IN clause
	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names))
	for i := range names {
		placeholders[i] = "?"
		args[i] = names[i]
	}

	query := fmt.Sprintf(`
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name IN (%s) AND deleted_at IS NULL
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *authRepository) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = ? AND deleted_at IS NULL
	`

	role := &domain.Role{}
	var deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if deletedAt.Valid {
		role.DeletedAt = &deletedAt.Time
	}

	return role, nil
}

func (r *authRepository) AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing roles
	deleteQuery := `DELETE FROM user_roles WHERE user_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, userID)
	if err != nil {
		return err
	}

	// Insert new roles
	insertQuery := `
		INSERT INTO user_roles (id, user_id, role_id, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	for _, roleID := range roleIDs {
		_, err = tx.ExecContext(ctx, insertQuery,
			uuid.New().String(),
			userID,
			roleID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *authRepository) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.deleted_at IS NULL AND ur.deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *authRepository) CreateUserToken(ctx context.Context, token *domain.UserToken) (*domain.UserToken, error) {
	query := `
		INSERT INTO user_tokens (id, user_id, token, ability, expired_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())`

	abilityJSON, err := json.Marshal(token.Ability)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ability: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		abilityJSON,
		token.ExpiredAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user token: %w", err)
	}

	return token, nil
}

func (r *authRepository) GetUserTokenByID(ctx context.Context, tokenID string) (*domain.UserToken, error) {
	query := `
		SELECT id, user_id, token, ability, expired_at, created_at, updated_at, deleted_at
		FROM user_tokens
		WHERE id = ? AND deleted_at IS NULL`

	userToken := &domain.UserToken{}
	var abilityJSON []byte
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tokenID).Scan(
		&userToken.ID,
		&userToken.UserID,
		&userToken.Token,
		&abilityJSON,
		&userToken.ExpiredAt,
		&userToken.CreatedAt,
		&userToken.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if deletedAt.Valid {
		userToken.DeletedAt = &deletedAt.Time
	}

	// Unmarshal ability JSON
	var ability []string
	if err := json.Unmarshal(abilityJSON, &ability); err != nil {
		return nil, err
	}
	userToken.Ability = ability

	return userToken, nil
}

func (r *authRepository) InvalidateUserToken(ctx context.Context, tokenID string) error {
	query := `
		UPDATE user_tokens
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("token not found")
	}

	return nil
}

func (r *authRepository) InvalidateUserTokens(ctx context.Context, userID string) error {
	query := `
		UPDATE user_tokens
		SET deleted_at = NOW()
		WHERE user_id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("no tokens found")
	}

	return nil
}
