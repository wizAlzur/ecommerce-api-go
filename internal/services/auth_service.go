package services

import (
	"context"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repositories"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (int64, error)
	Login(ctx context.Context, req *models.LoginRequest) (string, error)
	ValidateToken(tokenString string) (int64, error)
}

type authService struct {
	userRepo repositories.UserRepository
	jwtCfg   config.JWTConfig
}

func NewAuthService(userRepo repositories.UserRepository, jwtCfg config.JWTConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (int64, error) {
	hash, err := repositories.HashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("hash password error: %w", err)
	}
	return s.userRepo.Create(ctx, req.Email, hash)
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || !repositories.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", fmt.Errorf("invalid credentials")
	}

	exp := time.Now().Add(s.jwtCfg.Expires)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     exp.Unix(),
	})

	return token.SignedString([]byte(s.jwtCfg.Secret))
}

func (s *authService) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtCfg.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found or invalid")
	}

	return int64(userIDFloat), nil
}
