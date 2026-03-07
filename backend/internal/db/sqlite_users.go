package db

import (
	"database/sql"
	"fmt"

	"github.com/drews/basecoat/internal/models"
)

type sqliteUserRepo struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepo{db: db}
}

func (r *sqliteUserRepo) Create(u *models.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, username, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Email, u.PasswordHash, u.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *sqliteUserRepo) GetByID(id string) (*models.User, error) {
	return r.scanOne(r.db.QueryRow(`SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?`, id))
}

func (r *sqliteUserRepo) GetByEmail(email string) (*models.User, error) {
	return r.scanOne(r.db.QueryRow(`SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?`, email))
}

func (r *sqliteUserRepo) GetByUsername(username string) (*models.User, error) {
	return r.scanOne(r.db.QueryRow(`SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?`, username))
}

func (r *sqliteUserRepo) ListAll() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT id, username, email, password_hash, created_at FROM users`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	var out []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *sqliteUserRepo) scanOne(row *sql.Row) (*models.User, error) {
	u := &models.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return u, nil
}
