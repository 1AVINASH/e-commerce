package user

import (
	"database/sql"
	"fmt"
	"gotemplate/infra/postgres"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// GetUser fetches a single user by ID
func (r *UserRepository) GetUser(id int64) (*User, error) {
	var u User
	query := `SELECT id, name, email, password, session_token, role FROM users WHERE id = $1`
	err := postgres.DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("GetUser: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var u User
	query := `SELECT id, name, email, password, session_token, role FROM users WHERE email = $1`
	err := postgres.DB.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.SessionToken, &u.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("GetUser: %w", err)
	}
	return &u, nil
}

// ListUsers
func (r *UserRepository) GetUsers() ([]*User, error) {
	query := `SELECT id, name, email, role FROM users`
	rows, err := postgres.DB.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("GetUsers: %w", err)
	}

	users := []*User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// CreateUser inserts a new user into the DB
func (r *UserRepository) CreateUser(u *User) (*User, error) {
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`
	err := postgres.DB.QueryRow(query, *u.Name, *u.Email, *u.Password, *u.Role).Scan(&u.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}
	return u, nil
}
