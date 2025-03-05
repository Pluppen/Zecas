package services

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// GetByID returns a specific finding by ID
func (s *AuthService) GetBySessionToken(sessionToken string) (*models.Session, error) {
	var session models.Session
	result := s.db.Where("\"sessionToken\" = ?", sessionToken).First(&session)
	return &session, result.Error
}
