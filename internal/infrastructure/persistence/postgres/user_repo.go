package postgres

import (
	"database/sql"
	"errors"

	"banana-auction/internal/domain/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(u user.User) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO users (username, password_hash, name, role)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		u.Username, u.PasswordHash, u.Name, u.Role,
	).Scan(&id)
	if IsDuplicateKeyError(err) {
		return 0, errors.New("username already exists")
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepo) GetByUsername(username string) (user.User, error) {
	var u user.User
	err := r.db.QueryRow(`
		SELECT id, username, password_hash, name, role
		FROM users WHERE username = $1`, username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Name, &u.Role)
	if err == sql.ErrNoRows {
		return user.User{}, errors.New("user not found")
	}
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}

func (r *UserRepo) GetByID(id int) (user.User, error) {
	var u user.User
	err := r.db.QueryRow(`
		SELECT id, username, password_hash, name, role
		FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Name, &u.Role)
	if err == sql.ErrNoRows {
		return user.User{}, errors.New("user not found")
	}
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}
