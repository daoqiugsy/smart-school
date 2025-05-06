package model

import (
	"time"
)

// User 用户基础信息
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:100;not null" json:"-"`
	RealName  string    `gorm:"size:50" json:"real_name"`
	Email     string    `gorm:"size:100" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	UserType  int       `gorm:"not null;default:0" json:"user_type"` // 0:学生 1:教师 2:管理员
	Status    int       `gorm:"not null;default:1" json:"status"`    // 0:禁用 1:启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Student 学生信息
type Student struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	StudentID  string    `gorm:"size:20;not null;uniqueIndex" json:"student_id"` // 学号
	Grade      string    `gorm:"size:20" json:"grade"`                           // 年级
	Class      string    `gorm:"size:50" json:"class"`                           // 班级
	Major      string    `gorm:"size:50" json:"major"`                           // 专业
	Department string    `gorm:"size:50" json:"department"`                      // 院系
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
}

// TableName 指定表名
func (Student) TableName() string {
	return "students"
}

// Teacher 教师信息
type Teacher struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	TeacherID  string    `gorm:"size:20;not null;uniqueIndex" json:"teacher_id"` // 工号
	Title      string    `gorm:"size:50" json:"title"`                           // 职称
	Department string    `gorm:"size:50" json:"department"`                      // 院系
	Office     string    `gorm:"size:50" json:"office"`                          // 办公室
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
}

// TableName 指定表名
func (Teacher) TableName() string {
	return "teachers"
}
