package model

import (
	"time"
)

// Assignment 作业信息
type Assignment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CourseID    uint      `gorm:"not null;index" json:"course_id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      int       `gorm:"not null;default:1" json:"status"` // 0:已取消 1:进行中 2:已截止
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Course      Course    `gorm:"foreignKey:CourseID" json:"course"`
}

// TableName 指定表名
func (Assignment) TableName() string {
	return "assignments"
}

// StudentAssignment 学生作业提交记录
type StudentAssignment struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	StudentID    uint       `gorm:"not null;index" json:"student_id"`
	AssignmentID uint       `gorm:"not null;index" json:"assignment_id"`
	Content      string     `gorm:"type:text" json:"content"`
	Attachments  string     `gorm:"type:text" json:"attachments"`     // 附件路径，多个用逗号分隔
	Score        float64    `gorm:"default:null" json:"score"`        // 分数
	Comment      string     `gorm:"type:text" json:"comment"`         // 评语
	Status       int        `gorm:"not null;default:0" json:"status"` // 0:未提交 1:已提交 2:已批改
	SubmitTime   time.Time  `json:"submit_time"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Student      Student    `gorm:"foreignKey:StudentID" json:"student"`
	Assignment   Assignment `gorm:"foreignKey:AssignmentID" json:"assignment"`
}

// TableName 指定表名
func (StudentAssignment) TableName() string {
	return "student_assignments"
}

// Exam 考试信息
type Exam struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CourseID    uint      `gorm:"not null;index" json:"course_id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Location    string    `gorm:"size:100" json:"location"` // 考试地点
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ExamType    int       `gorm:"not null;default:0" json:"exam_type"` // 0:平时测验 1:期中考试 2:期末考试
	Status      int       `gorm:"not null;default:0" json:"status"`    // 0:未开始 1:进行中 2:已结束
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Course      Course    `gorm:"foreignKey:CourseID" json:"course"`
}

// TableName 指定表名
func (Exam) TableName() string {
	return "exams"
}
