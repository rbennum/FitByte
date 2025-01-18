package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/TimDebug/FitByte/dto"
	"github.com/TimDebug/FitByte/entity"
	"github.com/TimDebug/FitByte/helper"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{db: db}
}

func NewUserRepositoryInject(i do.Injector) (UserRepository, error) {
	db := do.MustInvoke[*pgxpool.Pool](i)
	return NewUserRepository(db), nil
}

func (r *UserRepository) GetProfile(ctx context.Context, id string) (*entity.User, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT email, name, image_uri FROM Users WHERE id = $1`,
		id,
	)

	var user entity.User
	err := row.Scan(&user.Email, &user.Name, &user.ImageUri)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetBatchOfProfiles(
	ctx context.Context,
	ids []string,
) ([]entity.User, error) {
	query := `
		SELECT id, email, name, image_uri
		FROM Users
		WHERE id = ANY($1::text[]);
	`
	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err) // Return wrapped errors
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.Id, &user.Email, &user.Name, &user.ImageUri)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) VerifyNewUser(ctx *gin.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM Users 
			WHERE email = $1
		);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		// the said user already exists
		return false, helper.ErrConflict
	}
	return true, nil
}

func (r *UserRepository) Login(ctx *gin.Context, body *dto.UserRequestPayload) ([]entity.User, error) {
	query := `
		SELECT id, email, password_hash
		FROM Users
		WHERE email = $1
	`
	rows, err := r.db.Query(ctx, query, body.Email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Email, &user.PasswordHash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Register(ctx *gin.Context, body *entity.User) (userId string, err error) {
	query := `
		INSERT INTO Users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING Users.id
	`
	createdAtUnix := time.Unix(body.CreatedAt, 0)
	row := r.db.QueryRow(ctx, query, body.Email, body.PasswordHash, createdAtUnix, createdAtUnix)
	err = row.Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, err
}
