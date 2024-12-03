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

	user.ID = uuid.New()
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

func (r *authRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password, name, phone, status, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Phone, &user.Status, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
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

func (r *authRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
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

func (r *authRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = ?
	`

	role := &domain.Role{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID, &role.Name, &role.Description,
		&role.CreatedAt, &role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *authRepository) AssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing roles
	_, err = tx.ExecContext(ctx, "DELETE FROM user_roles WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	// Insert new roles
	query := "INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)"
	for _, roleID := range roleIDs {
		_, err = tx.ExecContext(ctx, query, userID, roleID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *authRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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

func (r *authRepository) GetUserTokenByID(ctx context.Context, tokenID uuid.UUID) (*domain.UserToken, error) {
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

func (r *authRepository) InvalidateUserToken(ctx context.Context, tokenID uuid.UUID) error {
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

func (r *authRepository) InvalidateUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE user_tokens SET is_valid = false WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
