package repository

import (
	"smart-school/internal/model"

	"gorm.io/gorm"
)

// StudentRepository 学生数据访问接口
type StudentRepository interface {
	// Create 创建学生
	Create(student *model.Student) error
	// FindByID 根据ID查找学生
	FindByID(id uint) (*model.Student, error)
	// FindByUserID 根据用户ID查找学生
	FindByUserID(userID uint) (*model.Student, error)
	// FindByStudentID 根据学号查找学生
	FindByStudentID(studentID string) (*model.Student, error)
	// Update 更新学生信息
	Update(student *model.Student) error
	// Delete 删除学生
	Delete(id uint) error
}

// studentRepository 学生数据访问实现
type studentRepository struct {
	db *gorm.DB
}

// NewStudentRepository 创建学生仓库实例
func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

// Create 创建学生
func (r *studentRepository) Create(student *model.Student) error {
	return r.db.Create(student).Error
}

// FindByID 根据ID查找学生
func (r *studentRepository) FindByID(id uint) (*model.Student, error) {
	var student model.Student
	err := r.db.First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// FindByUserID 根据用户ID查找学生
func (r *studentRepository) FindByUserID(userID uint) (*model.Student, error) {
	var student model.Student
	err := r.db.Where("user_id = ?", userID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// FindByStudentID 根据学号查找学生
func (r *studentRepository) FindByStudentID(studentID string) (*model.Student, error) {
	var student model.Student
	err := r.db.Where("student_id = ?", studentID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// Update 更新学生信息
func (r *studentRepository) Update(student *model.Student) error {
	return r.db.Save(student).Error
}

// Delete 删除学生
func (r *studentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Student{}, id).Error
}
