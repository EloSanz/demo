package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/elosanz/demo/internal/user"
	"gorm.io/gorm"
)

// UserGORMRepository implements user.UserRepository using GORM.
type UserGORMRepository struct {
	db *gorm.DB
}

// NewUserGORMRepository constructs a UserGORMRepository backed by the given *gorm.DB.
func NewUserGORMRepository(db *gorm.DB) user.UserRepository {
	return &UserGORMRepository{db: db}
}

func (r *UserGORMRepository) FindAll(ctx context.Context, page int, size int) ([]user.User, int64, error) {
	var users []user.User
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting users: %w", err)
	}

	// Find page
	offset := (page - 1) * size
	if err := r.db.WithContext(ctx).Order("id asc").Offset(offset).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("querying users: %w", err)
	}

	if users == nil {
		users = []user.User{}
	}
	return users, total, nil
}

func (r *UserGORMRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by id %d: %w", id, err)
	}
	return &u, nil
}

func (r *UserGORMRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by email: %w", err)
	}
	return &u, nil
}

func (r *UserGORMRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("checking email existence: %w", err)
	}
	return count > 0, nil
}

func (r *UserGORMRepository) Save(ctx context.Context, u user.User) (user.User, error) {
	if err := r.db.WithContext(ctx).Save(&u).Error; err != nil {
		return user.User{}, fmt.Errorf("persisting user: %w", err)
	}
	return u, nil
}

func (r *UserGORMRepository) Update(ctx context.Context, u user.User) (user.User, error) {
	// Update with model to ensure CreatedAt is not zeroed out if not provided
	result := r.db.WithContext(ctx).Model(&u).Updates(user.User{
		Name:    u.Name,
		Email:   u.Email,
		Phone:   u.Phone,
		Website: u.Website,
		Points:  u.Points,
	})
	if result.Error != nil {
		return user.User{}, fmt.Errorf("updating user %d: %w", u.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return user.User{}, user.ErrNotFound
	}
	
	// Reload to get updated timestamps and all fields
	err := r.db.WithContext(ctx).First(&u, u.ID).Error
	return u, err
}

func (r *UserGORMRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&user.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("deleting user %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return user.ErrNotFound
	}
	return nil
}

func (r *UserGORMRepository) TransferPoints(ctx context.Context, fromID, toID int64, amount int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fromUser, toUser user.User

		// 1. Obtener y validar Usuario A (el que envía)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&fromUser, fromID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return user.ErrNotFound
			}
			return err
		}

		if fromUser.Points < amount {
			return user.ErrInsufficientPoints
		}

		// 2. Obtener Usuario B (el que recibe)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&toUser, toID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return user.ErrNotFound
			}
			return err
		}

		// 3. Ejecutar cambios
		fromUser.Points -= amount
		toUser.Points += amount

		if err := tx.Save(&fromUser).Error; err != nil {
			return err
		}

		if err := tx.Save(&toUser).Error; err != nil {
			return err
		}

		return nil // Si llegamos acá, GORM hace COMMIT automáticamente
	})
}
