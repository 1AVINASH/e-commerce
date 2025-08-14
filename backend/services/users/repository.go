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
	query := `SELECT id, name, email FROM users WHERE id = $1`
	err := postgres.DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email)
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
	query := `SELECT id, name, email FROM users`
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
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// CreateUser inserts a new user into the DB
func (r *UserRepository) CreateUser(u *User) (int64, error) {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := postgres.DB.QueryRow(query, u.Name, u.Email).Scan(&u.ID)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}
	return u.ID, nil
}
