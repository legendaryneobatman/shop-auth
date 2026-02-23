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
	query := fmt.Sprintf("INSERT INTO refresh_tokens (user_id,expires_at,ip_address,user_agent,revoked) values ($1, $2,$3, $4, $5)")
	row := r.db.QueryRow(
		query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.IPAddress,
		token.UserAgent,
		token.Revoked,
	)

	err := row.Scan()

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
	defer rows.Close()
	if rows.Err() != nil {
		logrus.Errorf("Error when rows for GetRefreshTokensByUserID %s", rows.Err().Error())
		return nil, rows.Err()
	}
	if err != nil {
		logrus.Errorf("Error when executing query for GetRefreshTokensByUserID %s", err.Error())
		return nil, err
	}

	var refreshTokens []models.RefreshToken
	for rows.Next() {
		var rf models.RefreshToken

		if err := rows.StructScan(&rf); err != nil {
			logrus.Fatalf("Error when scaning rows for GetRefreshTokensByUserID %s", err.Error())
			return nil, err
		}

		refreshTokens = append(refreshTokens, rf)
	}

	return refreshTokens, nil
}
func (r *Repository) RevokeRefreshToken(tokenID int) error {
	query := fmt.Sprintf("UPDATE refresh_tokens SET revoked=true WHERE $1")
	row := r.db.QueryRow(query, tokenID)
	err := row.Scan()
	if err != nil {
		logrus.Errorf("Error when RevokeRefreshToken %s", err.Error())
		return err
	}

	return nil
}
func (r *Repository) RevokeAllUserTokens(userID int) error {
	query := fmt.Sprintf("UPDATE refresh_tokens SET revoked=true WHERE user_id=$1")
	row := r.db.QueryRow(query, userID)
	err := row.Scan()

	if err != nil {
		logrus.Errorf("Error when RevokeAllUserTokens %s", err.Error())
	}

	return nil
}
