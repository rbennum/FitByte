package userRepository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/entity"
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

func (r *UserRepository) Create(ctx context.Context, user entity.User) (managerId string, err error) {
	query := `
		INSERT INTO manager (email, password)
		VALUES ($1, $2)
		RETURNING manager.managerid
	`
	fmt.Printf("email %s, password %s", user.Email.String, user.Password)
	row := r.db.QueryRow(ctx, query,
		user.Email.String, // Email yang unik
		user.Password,     // Kata sandi
	)

	err = row.Scan(&managerId)
	if err != nil {
		return "", err
	}

	return managerId, err
}
func (r *UserRepository) Update(ctx context.Context, user entity.User) error {
	query := `
		UPDATE manager
		SET name = $2, password = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		user.Id,          // UID
		user.Name.String, // Nama
		user.Password,    // Kata sandi
		user.UpdatedAt,   // Timestamp saat ini
	)
	return err
}

func (r *UserRepository) UpsertUser(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO manager (identitynumber, name, username, email, password, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			username = EXCLUDED.username,
			email = EXCLUDED.email,
			password = EXCLUDED.password,
			updated_at = EXCLUDED.updated_at
		WHERE 1=1
	`

	_, err := r.db.Exec(ctx, query,
		user.Id,              // UUID, jika kosong gunakan default UUID
		user.Name.String,     // Nama pengguna
		user.Username.String, // Username yang unik
		user.Email.String,    // Email yang unik
		user.Password,        // Kata sandi
		user.UpdatedAt,       // Timestamp saat ini
		user.CreatedAt,       // Timestamp saat dibuat
	)
	return err
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	query := `SELECT managerid, name, email FROM manager`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserbyEmail(ctx context.Context, email string) ([]entity.User, error) {
	// Menggunakan Query bukan Exec karena kita mengambil hasil dari SELECT
	query := `SELECT u.managerid, u.name, u.email, u.password FROM manager u WHERE u.email = $1`
	fmt.Printf("\nthefuck email %s\n", email)
	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		// Menyimpan data hasil query ke dalam struct user
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Memastikan tidak ada error saat iterasi
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM manager WHERE managerid = ?`
	_, err := r.db.Query(ctx, query, id)
	return err
}

func (r *UserRepository) GetProfile(ctx context.Context, id string) (*entity.GetProfile, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT email, name, userImageUri, companyName, companyImageUri FROM manager WHERE managerid = $1`,
		id,
	)

	var user entity.GetProfile
	err := row.Scan(&user.Email, &user.Name, &user.UserImageUri, &user.CompanyName, &user.CompanyImageUri)
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
		SELECT managerid, email, name, userImageUri, companyname, companyimageuri
		FROM manager
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
		err := rows.Scan(&user.ManagerId, &user.Email, &user.Name, &user.UserImageUri, &user.CompanyName, &user.CompanyImageUri)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id string, data *entity.GetProfile) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE manager SET 
			email= $1, 
			name=$2, 
			userImageUri=$3, 
			companyName=$4, 
			companyImageUri=$5
		WHERE managerid = $6`,
		data.Email,
		data.Name.String,
		data.UserImageUri.String,
		data.CompanyName.String,
		data.CompanyImageUri.String,
		id,
	)
	return err
}
