package model

import "time"

type Account struct {
	ID        uint
	Email     string
	Username  string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Online    bool
}

type AccountRegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
