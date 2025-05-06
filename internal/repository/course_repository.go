package repository

import (
	"gorm.io/gorm"
	"smart-school/internal/model"
)

// CourseRepository 课程数据访问接口
type CourseRepository interface {
	// Create 创建课程
	Create(course *model.Course) error
	// FindByID 根据ID查找课程
	FindByID(id uint) (*model.Course, error)
	// FindByCourseCode 根据课程代码查找课程
	FindByCourseCode(courseCode string) (*model.Course, error)
	// Update 更新课程信息
	Update(course *model.Course) error
	// Delete 删除课程
	Delete(id uint) error
	// GetCoursesByKeyword 根据关键词搜索课程
	GetCoursesByKeyword(keyword string) ([]model.Course, error)
	// GetStudentCourses 获取学生选课记录
	GetStudentCourses(studentID uint, semester string) ([]model.StudentCourse, error)
	// CreateStudentCourse 创建学生选课记录
	CreateStudentCourse(studentCourse *model.StudentCourse) error
	// GetCourseSchedules 获取课程安排
	GetCourseSchedules(courseID uint) ([]model.CourseSchedule, error)
	// CreateCourseSchedule 创建课程安排
	CreateCourseSchedule(schedule *model.CourseSchedule) error
}

// courseRepository 课程数据访问实现
type courseRepository struct {
	db *gorm.DB
}

// NewCourseRepository 创建课程仓库实例
func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

// Create 创建课程
func (r *courseRepository) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

// FindByID 根据ID查找课程
func (r *courseRepository) FindByID(id uint) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// FindByCourseCode 根据课程代码查找课程
func (r *courseRepository) FindByCourseCode(courseCode string) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("course_code = ?", courseCode).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// Update 更新课程信息
func (r *courseRepository) Update(course *model.Course) error {
	return r.db.Save(course).Error
}

// Delete 删除课程
func (r *courseRepository) Delete(id uint) error {
	return r.db.Delete(&model.Course{}, id).Error
}

// GetCoursesByKeyword 根据关键词搜索课程
func (r *courseRepository) GetCoursesByKeyword(keyword string) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Where("name LIKE ? OR description LIKE ? OR course_code LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(&courses).Error
	return courses, err
}

// GetStudentCourses 获取学生选课记录
func (r *courseRepository) GetStudentCourses(studentID uint, semester string) ([]model.StudentCourse, error) {
	var studentCourses []model.StudentCourse
	query := r.db.Where("student_id = ?", studentID).Preload("Course")

	if semester != "" {
		query = query.Joins("JOIN courses ON student_courses.course_id = courses.id").Where("courses.semester = ?", semester)
	}

	err := query.Find(&studentCourses).Error
	return studentCourses, err
}

// CreateStudentCourse 创建学生选课记录
func (r *courseRepository) CreateStudentCourse(studentCourse *model.StudentCourse) error {
	return r.db.Create(studentCourse).Error
}

// GetCourseSchedules 获取课程安排
func (r *courseRepository) GetCourseSchedules(courseID uint) ([]model.CourseSchedule, error) {
	var schedules []model.CourseSchedule
	err := r.db.Where("course_id = ?", courseID).Preload("Course").Preload("Teacher").Find(&schedules).Error
	return schedules, err
}

// CreateCourseSchedule 创建课程安排
func (r *courseRepository) CreateCourseSchedule(schedule *model.CourseSchedule) error {
	return r.db.Create(schedule).Error
}
