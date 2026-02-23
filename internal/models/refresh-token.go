package models

type RefreshToken struct {
	ID        int     `json:"id" db:"id"`
	UserID    int     `json:"user_id" db:"user_id"`
	TokenHash string  `json:"token_hash" db:"token_hash"`
	ExpiresAt string  `json:"expires_at" db:"expires_at"`
	IPAddress *string `json:"ip_address" db:"ip_address"`
	UserAgent *string `json:"user_agent" db:"user_agent"`
	Revoked   *bool   `json:"revoked" db:"revoked"`
}
