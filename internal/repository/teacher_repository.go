package repository

import (
	"smart-school/internal/model"

	"gorm.io/gorm"
)

// TeacherRepository 教师数据访问接口
type TeacherRepository interface {
	// Create 创建教师
	Create(teacher *model.Teacher) error
	// FindByID 根据ID查找教师
	FindByID(id uint) (*model.Teacher, error)
	// FindByUserID 根据用户ID查找教师
	FindByUserID(userID uint) (*model.Teacher, error)
	// FindByTeacherID 根据工号查找教师
	FindByTeacherID(teacherID string) (*model.Teacher, error)
	// Update 更新教师信息
	Update(teacher *model.Teacher) error
	// Delete 删除教师
	Delete(id uint) error
}

// teacherRepository 教师数据访问实现
type teacherRepository struct {
	db *gorm.DB
}

// NewTeacherRepository 创建教师仓库实例
func NewTeacherRepository(db *gorm.DB) TeacherRepository {
	return &teacherRepository{db: db}
}

// Create 创建教师
func (r *teacherRepository) Create(teacher *model.Teacher) error {
	return r.db.Create(teacher).Error
}

// FindByID 根据ID查找教师
func (r *teacherRepository) FindByID(id uint) (*model.Teacher, error) {
	var teacher model.Teacher
	err := r.db.First(&teacher, id).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// FindByUserID 根据用户ID查找教师
func (r *teacherRepository) FindByUserID(userID uint) (*model.Teacher, error) {
	var teacher model.Teacher
	err := r.db.Where("user_id = ?", userID).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// FindByTeacherID 根据工号查找教师
func (r *teacherRepository) FindByTeacherID(teacherID string) (*model.Teacher, error) {
	var teacher model.Teacher
	err := r.db.Where("teacher_id = ?", teacherID).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// Update 更新教师信息
func (r *teacherRepository) Update(teacher *model.Teacher) error {
	return r.db.Save(teacher).Error
}

// Delete 删除教师
func (r *teacherRepository) Delete(id uint) error {
	return r.db.Delete(&model.Teacher{}, id).Error
}
