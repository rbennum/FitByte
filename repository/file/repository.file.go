package fileRepository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/samber/do/v2"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) FileRepository {
	return FileRepository{db: db}
}

func NewInject(i do.Injector) (FileRepository, error) {
	_db := do.MustInvoke[*pgxpool.Pool](i)
	return NewFileRepository(_db), nil
}

func (r *FileRepository) Create(ctx context.Context, e entity.File) error {
	query := `
		INSERT INTO users (filename, fileuri)
		VALUES ($1, $2)
	`
	_, err := r.db.Exec(ctx, query,
		e.FileName, // Name
		e.FileURI,  // URI
	)
	return err
}

func (r *FileRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	query := `SELECT identitynumber, name, email FROM users`
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

func (r *FileRepository) GetFile(ctx context.Context, uri string) ([]entity.File, error) {
	// Menggunakan Query bukan Exec karena kita mengambil hasil dari SELECT
	query := `SELECT f.fileid, f.filename, f.fileuri FROM file f WHERE fileuri = $1`
	rows, err := r.db.Query(ctx, query, uri)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []entity.File
	for rows.Next() {
		var file entity.File
		// Menyimpan data hasil query ke dalam struct user
		if err := rows.Scan(&file.FileId, &file.FileName, &file.FileURI); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	// Memastikan tidak ada error saat iterasi
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *FileRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM file WHERE file.fileid = ?`
	_, err := r.db.Query(ctx, query, id)
	return err
}
