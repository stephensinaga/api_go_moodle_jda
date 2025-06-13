package model

import (
	"time"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}	

type ResponseData struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsExists  bool      `json:"isExists"`
	Token     string    `json:"token"`
}

type ReturnResponse struct {
	Data    ResponseData `json:"data"`
	Message string       `json:"message"`
	Status  bool         `json:"status"`
}
