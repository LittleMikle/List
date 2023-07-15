package service

import (
	"errors"
	"fmt"
	todo "github.com/LittleMikle/ToDo_List"
	"github.com/LittleMikle/ToDo_List/pkg/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/sha3"
	"time"
)

const (
	salt      = "sdg2klASs1Dlkjd0fmZd34dfgASc"
	signInKey = "SDfjk%0ncvk347dfs7sd72as!sd&"
	tokenTTL  = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	ExpiresAt int64
	IssuedAt  int64
	UserId    int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", fmt.Errorf("failed to generate token:%w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		UserId:    user.Id,
	})
	return token.SignedString([]byte(signInKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signInKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed with Parse with claims: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha3.New256()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
