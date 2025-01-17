package userRepository

import (
	"context"
	"fmt"

	"github.com/TimDebug/FitByte/entity"
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

func (r *UserRepository) GetProfile(ctx context.Context, id string) (*entity.GetProfile, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT email, name, image_uri FROM Users WHERE id = $1`,
		id,
	)

	var user entity.GetProfile
	err := row.Scan(&user.Email, &user.Name, &user.ImageUri)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetBatchOfProfiles(
	ctx context.Context,
	ids []string,
) ([]entity.GetProfile, error) {
	query := `
		SELECT managerid, email, name, image_uri
		FROM Users
		WHERE managerid = ANY($1::text[]);
	`
	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err) // Return wrapped errors
	}
	defer rows.Close()

	var users []entity.GetProfile
	for rows.Next() {
		var user entity.GetProfile
		err := rows.Scan(&user.Id, &user.Email, &user.Name, &user.ImageUri)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
