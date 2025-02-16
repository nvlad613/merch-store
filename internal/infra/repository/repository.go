package repository

import (
	"github.com/uptrace/bun"
)

type RepositoryImpl struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *RepositoryImpl {
	return &RepositoryImpl{db}
}
