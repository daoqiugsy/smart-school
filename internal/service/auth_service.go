package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"smart-school/internal/model"
	"smart-school/internal/repository"
	"time"
)

// 定义错误常量
var (
	ErrUserExists         = errors.New("用户名已存在")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("用户已被禁用")
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	RealName  string `json:"real_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserType  *int   `json:"user_type" binding:"required"` // 0:学生 1:教师 2:管理员
	StudentID string `json:"student_id"`                   // 学生学号，学生注册时必填
	TeacherID string `json:"teacher_id"`                   // 教师工号，教师注册时必填
	// 学生特有字段
	Grade      string `json:"grade"`
	Class      string `json:"class"`
	Major      string `json:"major"`
	Department string `json:"department"`
	// 教师特有字段
	Title  string `json:"title"`
	Office string `json:"office"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	RealName  string    `json:"real_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	UserType  int       `json:"user_type"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthService 认证服务接口
type AuthService interface {
	// Register 用户注册
	Register(req RegisterRequest) error
	// Login 用户登录
	Login(req LoginRequest) (*UserResponse, error)
}

// authService 认证服务实现
type authService struct {
	userRepo    repository.UserRepository
	studentRepo repository.StudentRepository
	teacherRepo repository.TeacherRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, studentRepo repository.StudentRepository, teacherRepo repository.TeacherRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		studentRepo: studentRepo,
		teacherRepo: teacherRepo,
	}
}

// Register 用户注册
func (s *authService) Register(req RegisterRequest) error {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(req.Username)
	if err == nil && existingUser != nil {
		return ErrUserExists
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	now := time.Now()
	user := &model.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		RealName:  req.RealName,
		Email:     req.Email,
		Phone:     req.Phone,
		UserType:  *req.UserType,
		Status:    1, // 默认启用
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 保存用户
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// 根据用户类型创建对应的角色信息
	switch *req.UserType {
	case 0: // 学生
		student := &model.Student{
			UserID:     user.ID,
			StudentID:  req.StudentID,
			Grade:      req.Grade,
			Class:      req.Class,
			Major:      req.Major,
			Department: req.Department,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		return s.studentRepo.Create(student)

	case 1: // 教师
		teacher := &model.Teacher{
			UserID:     user.ID,
			TeacherID:  req.TeacherID,
			Title:      req.Title,
			Department: req.Department,
			Office:     req.Office,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		return s.teacherRepo.Create(teacher)
	}

	return nil
}

// Login 用户登录
func (s *authService) Login(req LoginRequest) (*UserResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, ErrUserDisabled
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 返回用户信息
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		RealName:  user.RealName,
		Email:     user.Email,
		Phone:     user.Phone,
		UserType:  user.UserType,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
