package user

import (
	"context"
	"strings"
	"time"

	"boilerplate/internal/user/model"
	"boilerplate/pkg/jwt"
	"boilerplate/pkg/logger"
	"boilerplate/pkg/redis"
	userErr "boilerplate/shared/errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db          *gorm.DB
	jwtSecret   string
	logger      logger.Logger
	redisClient *redis.RedisClient
}

func NewUserService(db *gorm.DB, jwtSecret string, logger logger.Logger, redisClient *redis.RedisClient) *UserService {
	if db == nil {
		panic("database connection is required")
	}
	if logger == nil {
		panic("logger is required")
	}
	if redisClient == nil {
		panic("redis client is required")
	}
	if jwtSecret == "" {
		panic("jwt secret is required")
	}

	return &UserService{
		db:          db,
		jwtSecret:   jwtSecret,
		logger:      logger,
		redisClient: redisClient,
	}
}

func (s *UserService) Register(input model.RegisterInput) (*model.User, error) {
	var existingUser model.User
	if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": userErr.ErrEmailAlreadyRegistered.Error(),
		}).Error("Email sudah terdaftar")
		return nil, userErr.ErrEmailAlreadyRegistered
	}

	user, err := model.NewUser(input)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Error("Gagal membuat user baru")
		return nil, err
	}

	if err := s.db.Create(user).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": err.Error(),
		}).Error("Gagal menyimpan user ke database")
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(input model.LoginInput) (string, error) {
	var user model.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": userErr.ErrInvalidCredentials.Error(),
		}).Error("Kredensial login tidak valid")
		return "", userErr.ErrInvalidCredentials
	}

	if err := user.CheckPassword(input.Password); err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": input.Email,
			"error": userErr.ErrInvalidCredentials.Error(),
		}).Error("Password tidak valid")
		return "", userErr.ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("Gagal generate token")
		return "", err
	}

	// Simpan token di Redis
	ctx := context.Background()
	if err := s.redisClient.SetToken(ctx, user.ID, token, 24*time.Hour); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("Gagal menyimpan token di Redis")
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetStoredToken(ctx context.Context, userID uint) (string, error) {
	if s.redisClient == nil {
		return "", userErr.ErrInvalidCredentials
	}
	return s.redisClient.GetToken(ctx, userID)
}

func (s *UserService) Logout(userID uint) error {
	ctx := context.Background()
	if err := s.redisClient.DeleteToken(ctx, userID); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Gagal menghapus token dari Redis")
		return err
	}
	return nil
}

func (s *UserService) UpdateProfile(userID uint, input model.UpdateProfileInput) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Update fields
	user.Name = strings.TrimSpace(input.Name)
	if input.Email != "" {
		user.Email = strings.TrimSpace(input.Email)
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(input.Password)), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) DeleteAccount(userID uint) error {
	return s.db.Delete(&model.User{}, userID).Error
}

// GetAllUsers mengambil semua user tanpa password
func (s *UserService) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := s.db.Find(&users).Error; err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Gagal mengambil daftar user")
		return nil, err
	}

	// Hapus password dari setiap user
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}
