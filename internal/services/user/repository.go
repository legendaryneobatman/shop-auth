package user

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (h *Handler) CreateUser()  {}
func (h *Handler) GetByIdUser() {}
func (h *Handler) GetListUser() {}
func (h *Handler) EditUser()    {}
func (h *Handler) DeleteUser()  {}
