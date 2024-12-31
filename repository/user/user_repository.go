package userRepository

import (
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(music entity.User, txOpt *TxOptions) error
	FindByID(id string, txOpt *TxOptions) (bool, error)
	FindByEmail(email string, txOpt *TxOptions) (entity.User, error)
	Update(user entity.User, txOpt *TxOptions) (entity.User, error)
	DeleteByID(ID string, txOpt *TxOptions) error
	IsEmailRegistered(email string, txOpt *TxOptions) error
	IsUsernameRegistered(username string, txOpt *TxOptions) error
}

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repository{db}
}

type TxOptions struct {
	Tx *gorm.DB
}

func (txOpt TxOptions) IsTransactionActive() bool {
	return txOpt.Tx.Statement != nil && txOpt.Tx.Statement.DB.ConnPool != nil
}

func (r repository) Create(user entity.User, txOpt *TxOptions) error {
	sql := r.db
	if txOpt != nil && txOpt.Tx != nil && txOpt.IsTransactionActive() {
		sql = txOpt.Tx
	}

	query := `
	INSERT INTO users (
		id,
		name,
		username,
		email,
		password,
		created_at,
		updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT DO NOTHING`

	err := sql.Exec(
		query,
		user.Id,
		user.Name,
		user.Username,
		user.Password,
		user.UpdatedAt,
		user.CreatedAt,
	).Error

	if err != nil {
		return err
	}
	return nil
}

func (r repository) FindByID(id string, txOpt *TxOptions) (bool, error) {
	sql := r.db
	if txOpt != nil && txOpt.Tx != nil && txOpt.IsTransactionActive() {
		sql = txOpt.Tx
	}

	var count int64
	query := `SELECT 1 FROM users WHERE id = ?`
	err := sql.Raw(query, id).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r repository) FindByEmail(email string, txOpt *TxOptions) (entity.User, error) {
	sql := r.db
	if txOpt != nil && txOpt.Tx != nil && txOpt.IsTransactionActive() {
		sql = txOpt.Tx
	}

	var user entity.User
	query := `SELECT * FROM users WHERE email = ?`
	err := sql.Raw(query, email).Scan(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r repository) Update(user entity.User, txOpt *TxOptions) (entity.User, error) {
	sql := r.db
	if txOpt != nil && txOpt.Tx != nil && txOpt.IsTransactionActive() {
		sql = txOpt.Tx
	}

	query := `
	UPDATE users 
	SET 
		name = ?, 
		username = ?,
		email = ?, 
		password = ?, 
	WHERE id = ?`
	err := sql.Exec(
		query,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Id,
	).Error

	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r repository) DeleteByID(ID string, txOpt *TxOptions) error {
	sql := r.db
	if txOpt != nil && txOpt.Tx != nil && txOpt.IsTransactionActive() {
		sql = txOpt.Tx
	}

	query := `DELETE FROM users WHERE id = ?`
	err := sql.Exec(query, ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r repository) IsEmailRegistered(email string, txOpt *TxOptions) error {
	sql := r.db
	if txOpt != nil && txOpt.Tx != nil && txOpt.IsTransactionActive() {
		sql = txOpt.Tx
	}

	var count int64
	query := `SELECT 1 FROM users WHERE email = ?`
	err := sql.Raw(query, email).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return helper.ErrorEmailRegistered
	}
	return nil
}

func (r repository) IsUsernameRegistered(username string, txOpt *TxOptions) error {
	sql := r.db
	if txOpt != nil && txOpt.Tx != nil && txOpt.IsTransactionActive() {
		sql = txOpt.Tx
	}

	var count int64
	query := `SELECT 1 FROM users WHERE username = ?`
	err := sql.Raw(query, username).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return helper.ErrorUsernameRegistered
	}
	return nil
}
