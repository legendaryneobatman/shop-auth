package auth

import (
	"fmt"
	"shop-auth/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(
		"SELECT id, name, username, password_hash, avatar_url, email, phone, role, is_active, created_at, updated_at FROM users WHERE username=$1",
	)
	err := r.db.Get(&user, query, username)
	if err != nil {
		logrus.Errorf("Error when GetUserByUsername %s", err.Error())
		return nil, err
	}

	return &user, nil
}
func (r *Repository) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(
		"SELECT id, name, username, password_hash, avatar_url, email, phone, role, is_active, created_at, updated_at FROM users WHERE id=$1",
	)
	err := r.db.Get(&user, query, userID)
	if err != nil {
		logrus.Errorf("Error when GetUserByID %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *Repository) SaveRefreshToken(token models.RefreshToken) error {
	query := fmt.Sprintf(`
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, ip_address, user_agent, revoked)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	_, err := r.db.Exec(
		query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.IPAddress,
		token.UserAgent,
		token.Revoked,
	)
	if err != nil {
		logrus.Errorf("Error when SaveRefreshToken %s", err.Error())
		return err
	}
	return nil
}
func (r *Repository) GetRefreshTokenByHash(hash string) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}

	query := fmt.Sprintf("SELECT id, user_id, token_hash, expires_at, ip_address, user_agent, revoked FROM refresh_tokens WHERE token_hash=$1")
	err := r.db.Get(refreshToken, query, hash)
	if err != nil {
		logrus.Errorf("Error when GetRefreshTokenByHash %s", err.Error())
		return nil, err
	}

	return refreshToken, nil
}
func (r *Repository) GetRefreshTokensByUserID(userID int) ([]models.RefreshToken, error) {
	query := fmt.Sprintf("SELECT id, user_id, token_hash, expires_at, ip_address, user_agent, revoked FROM refresh_tokens WHERE user_id=$1")
	rows, err := r.db.Queryx(query, userID)
	if err != nil {
		logrus.Errorf("Error when executing query for GetRefreshTokensByUserID %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var refreshTokens []models.RefreshToken
	for rows.Next() {
		var rf models.RefreshToken

		if err := rows.StructScan(&rf); err != nil {
			logrus.Errorf("Error when scaning rows for GetRefreshTokensByUserID %s", err.Error())
			return nil, err
		}

		refreshTokens = append(refreshTokens, rf)
	}

	if rows.Err() != nil {
		logrus.Errorf("Error when iterating rows for GetRefreshTokensByUserID %s", rows.Err().Error())
		return nil, rows.Err()
	}

	return refreshTokens, nil
}
func (r *Repository) RevokeRefreshToken(tokenID int) error {
	query := fmt.Sprintf("UPDATE refresh_tokens SET revoked=true WHERE id=$1")
	_, err := r.db.Exec(query, tokenID)
	if err != nil {
		logrus.Errorf("Error when RevokeRefreshToken %s", err.Error())
		return err
	}

	return nil
}
func (r *Repository) RevokeAllUserTokens(userID int) error {
	query := fmt.Sprintf("UPDATE refresh_tokens SET revoked=true WHERE user_id=$1")
	_, err := r.db.Exec(query, userID)
	if err != nil {
		logrus.Errorf("Error when RevokeAllUserTokens %s", err.Error())
		return err
	}
	return nil
}

func (r *Repository) CreateUser(user models.User) (int, error) {
	query := fmt.Sprintf(`
		INSERT INTO users (name, username, password_hash, avatar_url, email, phone, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id
	`)

	var id int
	if err := r.db.QueryRow(
		query,
		user.Name,
		user.Username,
		user.Password,
		user.AvatarURL,
		user.Email,
		user.Phone,
		user.Role,
		user.IsActive,
	).Scan(&id); err != nil {
		logrus.Errorf("Error when CreateUser %s", err.Error())
		return 0, err
	}

	return id, nil
}
