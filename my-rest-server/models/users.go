package models

import (
	"time"

	"golang.org/x/oauth2"
)

type User struct {
	oauth2.Token
	UserId     int
	Email      string
	FullName   string
	CreatedAt  time.Time
	CoverPhoto string
}

func GetUserTokens(userId int) (oauth2.Token, error) {
	t := oauth2.Token{}
	row := DB.QueryRow("SELECT access_token, refresh_token, expires_in, expiry, token_type from users WHERE user_id=?", userId)
	err := row.Scan(&t.AccessToken, &t.RefreshToken, &t.ExpiresIn, &t.Expiry, &t.TokenType)

	return t, err
}

// email, full_name, cover_photo
func ListUserById(userId int) (User, error) {
	var user User

	rows, err := DB.Query("SELECT email, full_name, cover_photo from users WHERE user_id=?", userId)
	if err != nil {
		return User{}, err
	}

	for rows.Next() {
		rows.Scan(&user.Email, &user.FullName, &user.CoverPhoto)
	}

	return user, nil
}

// returns {Email, FullName, CreatedAt}[]
func ListUsers() ([]User, error) {
	var users []User

	rows, err := DB.Query("SELECT user_id, email, full_name, cover_photo, created_at FROM users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		rows.Scan(&user.UserId, &user.Email, &user.FullName, &user.CoverPhoto, &user.CreatedAt)

		users = append(users, user)
	}

	return users, nil
}

func AddUser(u User) error {
	_, err := DB.Exec("INSERT INTO users(email, full_name, cover_photo, access_token, refresh_token, expires_in, expiry) VALUES (?, ?, ?, ?, ?, ?, ?)", u.Email, u.FullName, u.CoverPhoto, u.AccessToken, u.RefreshToken, u.ExpiresIn, u.Expiry)
	return err
}
