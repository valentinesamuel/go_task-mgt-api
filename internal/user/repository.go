package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, errors.New("user is empty")
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is empty")
	}

	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil || user.ID == 0 {
		return nil, errors.New("invalid user or user id")
	}

	var existingUser models.User
	err := r.db.WithContext(ctx).First(&existingUser, user.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with id %d", user.ID)
		}

	}

	err = r.db.WithContext(ctx).Model(&existingUser).Updates(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}

	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = r.db.WithContext(ctx).Delete(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &user, nil
}
