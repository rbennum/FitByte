package db_trx_repository

import (
	"gorm.io/gorm"
)

type DBTrxRepository interface {
	Begin() (*gorm.DB, error)
	Rollback(tx *gorm.DB)
	Commit(tx *gorm.DB)
}

type repository struct {
	db *gorm.DB
}

func NewDBTrxRepository(db *gorm.DB) DBTrxRepository {
	return &repository{db}
}

func (t *repository) Begin() (*gorm.DB, error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (t *repository) Rollback(tx *gorm.DB) {
	tx.Rollback()
}

func (t *repository) Commit(tx *gorm.DB) {
	tx.Commit()
}
