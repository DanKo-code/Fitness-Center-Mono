package postgres

import (
	"User/internal/dtos"
	customErrors "User/internal/errors"
	"User/internal/models"
	"User/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

var (
	clientRole = "client"
	adminRole  = "admin"
)

func (userRep *UserRepository) UpdateUser(ctx context.Context, cmd *dtos.UpdateUserCommand) error {
	setFields := map[string]interface{}{}

	if cmd.Name != "" {
		setFields["name"] = cmd.Name
	}
	if cmd.Email != "" {
		setFields["email"] = cmd.Email
	}
	if cmd.Role != "" {
		setFields["role"] = cmd.Role
	}
	if cmd.Photo != "" {
		setFields["photo"] = cmd.Photo
	}
	setFields["updated_time"] = cmd.UpdatedTime

	if len(setFields) == 0 {
		logger.InfoLogger.Printf("No fields to update for user Id: %v", cmd.Id)
		return nil
	}

	query := `UPDATE "user" SET `

	var params []interface{}
	i := 1
	for field, value := range setFields {
		if i > 1 {
			query += ", "
		}

		query += fmt.Sprintf(`%s = $%d`, field, i)
		params = append(params, value)
		i++
	}
	query += fmt.Sprintf(` WHERE id = $%d`, i)
	params = append(params, cmd.Id)

	_, err := userRep.db.ExecContext(ctx, query, params...)
	if err != nil {
		logger.ErrorLogger.Printf("Error UpdateUser: %v", err)
		return err
	}

	return nil
}

func (userRep *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := userRep.db.NamedExecContext(ctx, `
		INSERT INTO "user" (id, name, email, role, password_hash, photo, created_time, updated_time)
		VALUES (:id, :name, :email, :role, :password_hash, :photo, :created_time, :updated_time)`, *user)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateUser: %v", err)
		return nil, err
	}

	return user, nil
}

func (userRep *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := userRep.db.GetContext(ctx, user, `SELECT id, name, email, role, password_hash, photo, created_time, updated_time FROM "user" WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetUserById: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (userRep *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := userRep.db.GetContext(ctx, user, `SELECT id, name, email, role, password_hash, photo, created_time, updated_time FROM "user" WHERE email = $1`, email)
	if err != nil {
		logger.ErrorLogger.Printf("Error GetUserByEmail: %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (userRep *UserRepository) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	_, err := userRep.db.ExecContext(ctx, `DELETE FROM "user" WHERE id = $1`, id)
	if err != nil {
		logger.ErrorLogger.Printf("Error DeleteUserById: %v", err)
		return err
	}

	return nil
}

func (userRep *UserRepository) GetUsersByIds(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	if len(ids) == 0 {
		return []*models.User{}, nil
	}

	query := `
		SELECT id, email, role, photo, name, created_time, updated_time
		FROM "user"
		WHERE id = ANY($1)
	`

	var users []*models.User

	err := userRep.db.SelectContext(ctx, &users, query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to get users by ids: %w", err)
	}

	return users, nil
}

func (userRep *UserRepository) GetClients(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, email, role, photo, name, created_time, updated_time
		FROM "user"
		WHERE role = $1
	`

	var users []*models.User

	err := userRep.db.SelectContext(ctx, &users, query, clientRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	return users, nil
}

func (userRep *UserRepository) GetAdmins(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, email, role, photo, name, created_time, updated_time
		FROM "user"
		WHERE role = $1
	`

	var users []*models.User

	err := userRep.db.SelectContext(ctx, &users, query, adminRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get admins: %w", err)
	}

	return users, nil
}
