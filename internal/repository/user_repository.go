package repository

import (
	"smart-school/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// Create 创建用户
	Create(user *model.User) error
	// FindByID 根据ID查找用户
	FindByID(id uint) (*model.User, error)
	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*model.User, error)
	// Update 更新用户信息
	Update(user *model.User) error
	// Delete 删除用户
	Delete(id uint) error
}

// userRepository 用户数据访问实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
