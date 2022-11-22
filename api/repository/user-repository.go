package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"test/api/entity"

	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Register(ctx context.Context, tx *sql.Tx, user entity.Users) error
	GetLastInsertUser(ctx context.Context, tx *sql.Tx) (entity.Users, error)
	CheckUserByEmail(ctx context.Context, tx *sql.Tx, email string) error
	UpdateVerifiedAt(ctx context.Context, tx *sql.Tx, timeVerifiedAt time.Time, email string) error
}

type userRepository struct {
	// handlers.UserController
}

func NewUserRepository() UserRepository {
	return &userRepository{}

}

// Register is a function to register a new user to the database
func (repository *userRepository) Register(ctx context.Context, tx *sql.Tx, user entity.Users) error {
	var id int
	var email, username string
	var emailArr, usernameArr []string
	rowsCheck, err := tx.QueryContext(ctx, "SELECT email, username FROM users")

	if err != nil {
		return err
	}

	for rowsCheck.Next() {
		rowsCheck.Scan(&email, &username)
		emailArr = append(emailArr, email)
		usernameArr = append(usernameArr, username)
	}

	for _, value := range usernameArr {
		if value == user.Username {
			return fmt.Errorf("username has been registered")
		}
	}

	for _, value := range emailArr {
		if value == user.Email {
			return fmt.Errorf("email has been registered")
		}
	}

	temp, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	user.Password = string(temp)

	_, err = tx.ExecContext(ctx, "INSERT INTO users (name, username, email, password, role, email_verification, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", user.Name, user.Username, user.Email, user.Password, user.Role, user.EmailVerification, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return err
	}

	rows := tx.QueryRowContext(ctx, "SELECT id FROM users WHERE username = ?", user.Username)

	rows.Scan(&id)

	_, err = tx.ExecContext(ctx, "INSERT INTO user_details (user_id, phone, gender, type_of_disability, birthdate) VALUES (?, ?, ?, ?, ?)", id, user.Phone, user.Gender, user.DisabilityType, user.Birthdate)

	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) CheckUserByEmail(ctx context.Context, tx *sql.Tx, email string) error {
	query := "SELECT * FROM users WHERE email = ?"
	queryContext, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		return err
	}
	defer func(queryContext *sql.Rows) {
		err := queryContext.Close()
		if err != nil {
			return
		}
	}(queryContext)

	if queryContext.Next() {
		return nil
	}

	return errors.New("the user with the email was not found")
}
func (repository *userRepository) UpdateVerifiedAt(ctx context.Context, tx *sql.Tx, timeVerifiedAt time.Time, email string) error {
	query := "UPDATE users SET email_verification = ? WHERE email = ?"
	_, err := tx.ExecContext(ctx, query, timeVerifiedAt, email)
	if err != nil {
		return err
	}

	return nil
}

// GetLastInsertUser is a function to get the last inserted user from the database
func (repository *userRepository) GetLastInsertUser(ctx context.Context, tx *sql.Tx) (entity.Users, error) {

	var user entity.Users

	rows := tx.QueryRowContext(ctx, "SELECT users.id, users.name, users.username, users.email, users.password, users.role, user_details.phone, user_details.gender, user_details.type_of_disability, user_details.birthdate, users.email_verification, users.created_at, users.updated_at FROM users INNER JOIN user_details ON user_details.user_id = users.id WHERE users.id = (SELECT MAX(id) FROM users)")

	rows.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Role, &user.Phone, &user.Gender, &user.DisabilityType, &user.Birthdate, &user.EmailVerification, &user.CreatedAt, &user.UpdatedAt)
	return user, nil
}
