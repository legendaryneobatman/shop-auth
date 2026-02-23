package models

import "time"

type User struct {
	ID        int       `json:"-" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Username  string    `json:"username" db:"username" binding:"required"`
	Password  string    `json:"password" db:"password_hash" binding:"required"`
	AvatarURL *string   `json:"avatar_url,omitempty" db:"avatar_url"` // Опциональное поле
	Email     *string   `json:"email" db:"email"`
	Phone     *string   `json:"phone,omitempty" db:"phone"` // Опциональное поле
	Role      *string   `json:"role" db:"role"`             // Например: "user", "admin", "moderator"
	IsActive  bool      `json:"is_active" db:"is_active"`   // Для блокировки/активации
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
