package service

import (
	"avito_go/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	salt       = "kemfoefo,s1313l;akpod1012e"
	tokentTTL  = 12 * time.Hour
	signingKey = "kwjf203iraslf"
)

type AuthService struct {
	repository repository.Authorization
}

func NewAuthService(repository repository.Authorization) *AuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
	Coins  int `json:"coins"`
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repository.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{ExpiresAt: time.Now().Add(tokentTTL).Unix(),
			IssuedAt: time.Now().Unix()}, user.Id, user.Coins})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(token string) (int, int, error) {
	t, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, 0, err
	}

	claims, ok := t.Claims.(*tokenClaims)

	if !ok {
		return 0, 0, errors.New("token claims in not of type")
	}

	return claims.UserId, claims.Coins, nil
}
