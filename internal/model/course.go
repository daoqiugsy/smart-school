package model

import (
	"time"
)

// Course 课程信息
type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CourseCode  string    `gorm:"size:20;not null;uniqueIndex" json:"course_code"` // 课程代码
	Name        string    `gorm:"size:100;not null" json:"name"`                   // 课程名称
	Description string    `gorm:"type:text" json:"description"`                    // 课程描述
	Credit      float64   `gorm:"not null;default:0" json:"credit"`                // 学分
	Semester    string    `gorm:"size:20;not null" json:"semester"`                // 学期
	Status      int       `gorm:"not null;default:1" json:"status"`                // 0:未开课 1:已开课 2:已结课
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}

// CourseSchedule 课程安排
type CourseSchedule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"not null;index" json:"course_id"`
	TeacherID uint      `gorm:"not null;index" json:"teacher_id"`
	Classroom string    `gorm:"size:50" json:"classroom"`           // 教室
	Building  string    `gorm:"size:50" json:"building"`            // 教学楼
	Weekday   int       `gorm:"not null" json:"weekday"`            // 星期几 (1-7)
	StartWeek int       `gorm:"not null" json:"start_week"`         // 开始周次
	EndWeek   int       `gorm:"not null" json:"end_week"`           // 结束周次
	StartTime string    `gorm:"size:10;not null" json:"start_time"` // 开始时间 (HH:MM)
	EndTime   string    `gorm:"size:10;not null" json:"end_time"`   // 结束时间 (HH:MM)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Course    Course    `gorm:"foreignKey:CourseID" json:"course"`
	Teacher   Teacher   `gorm:"foreignKey:TeacherID" json:"teacher"`
}

// TableName 指定表名
func (CourseSchedule) TableName() string {
	return "course_schedules"
}

// StudentCourse 学生选课记录
type StudentCourse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID uint      `gorm:"not null;index" json:"student_id"`
	CourseID  uint      `gorm:"not null;index" json:"course_id"`
	Score     float64   `gorm:"default:null" json:"score"`        // 成绩
	Status    int       `gorm:"not null;default:1" json:"status"` // 0:退选 1:已选 2:已修完
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Student   Student   `gorm:"foreignKey:StudentID" json:"student"`
	Course    Course    `gorm:"foreignKey:CourseID" json:"course"`
}

// TableName 指定表名
func (StudentCourse) TableName() string {
	return "student_courses"
}
