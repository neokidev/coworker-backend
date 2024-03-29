package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ot07/coworker-backend/util"
)

var (
	ErrExpiredToken = errors.New("token has expired")
)

// Token contains the token and expired at
type Token struct {
	ID        uuid.UUID `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewToken creates a new token with a specific duration
func NewToken(duration time.Duration) *Token {
	tokenID := util.RandomUUID()

	token := &Token{
		ID:        tokenID,
		ExpiredAt: time.Now().Add(duration),
	}
	return token
}

// Valid checks if the token is valid or not
func (token *Token) Valid() error {
	if time.Now().After(token.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
