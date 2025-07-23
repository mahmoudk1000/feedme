package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreateAt  time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"apikey"`
}

func DatabaseUserToUser(u database.User) User {
	return User{
		Id:        u.ID,
		CreateAt:  u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name,
		ApiKey:    u.ApiKey,
	}
}
