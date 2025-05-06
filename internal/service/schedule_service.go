package service

import (
	"encoding/csv"
	"errors"
	"io"
	"smart-school/internal/model"
	"smart-school/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ScheduleService 课程表服务接口
type ScheduleService interface {
	// ImportFromExcel 从Excel导入课程表
	ImportFromExcel(studentID uint, file io.Reader) error
	// ImportFromCSV 从CSV导入课程表
	ImportFromCSV(studentID uint, file io.Reader) error
	// ImportFromAPI 从教务系统API导入课程表
	ImportFromAPI(studentID uint, apiURL, username, password string) error
	// GetStudentSchedule 获取学生课程表
	GetStudentSchedule(studentID uint) ([]model.CourseSchedule, error)
}

type scheduleService struct {
	studentRepo repository.StudentRepository
	courseRepo  repository.CourseRepository
}

// NewScheduleService 创建课程表服务实例
func NewScheduleService(studentRepo repository.StudentRepository, courseRepo repository.CourseRepository) ScheduleService {
	return &scheduleService{
		studentRepo: studentRepo,
		courseRepo:  courseRepo,
	}
}

// ImportFromExcel 从Excel导入课程表
func (s *scheduleService) ImportFromExcel(studentID uint, file io.Reader) error {
	// 解析Excel文件
	xlsFile, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}
	defer xlsFile.Close()

	// 获取第一个工作表
	sheetName := xlsFile.GetSheetName(0)
	rows, err := xlsFile.GetRows(sheetName)
	if err != nil {
		return err
	}

	// 解析课程数据并保存
	// 这里需要根据实际Excel格式进行调整
	return s.processCourseData(studentID, rows)
}

// ImportFromCSV 从CSV导入课程表
func (s *scheduleService) ImportFromCSV(studentID uint, file io.Reader) error {
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// 解析课程数据并保存
	return s.processCourseData(studentID, rows)
}

// ImportFromAPI 从教务系统API导入课程表
func (s *scheduleService) ImportFromAPI(studentID uint, apiURL, username, password string) error {
	// 实现教务系统API对接
	// 1. 登录教务系统
	// 2. 获取课程表数据
	// 3. 解析数据并保存

	// 这里是示例代码，需要根据实际教务系统API进行实现
	return errors.New("教务系统API对接功能尚未实现")
}

// GetStudentSchedule 获取学生课程表
func (s *scheduleService) GetStudentSchedule(userID uint) ([]model.CourseSchedule, error) {
	// 根据用户ID获取学生信息
	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 获取学生选课记录
	courses, err := s.courseRepo.GetStudentCourses(student.ID, "")
	if err != nil {
		return nil, err
	}

	// 获取课程安排
	var schedules []model.CourseSchedule
	for _, course := range courses {
		courseSchedules, err := s.courseRepo.GetCourseSchedules(course.CourseID)
		if err != nil {
			continue
		}
		schedules = append(schedules, courseSchedules...)
	}

	return schedules, nil
}

// processCourseData 处理课程数据并保存
func (s *scheduleService) processCourseData(studentID uint, rows [][]string) error {
	// 解析课程数据
	// 这里需要根据实际数据格式进行调整

	// 示例实现
	for i, row := range rows {
		// 跳过表头
		if i == 0 {
			continue
		}

		// 确保行数据足够
		if len(row) < 7 {
			continue
		}

		// 解析课程信息
		courseCode := strings.TrimSpace(row[0])
		courseName := strings.TrimSpace(row[1])
		credit, _ := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		semester := strings.TrimSpace(row[3])

		// 创建或更新课程
		course, err := s.courseRepo.FindByCourseCode(courseCode)
		if err != nil || course == nil {
			// 创建新课程
			course = &model.Course{
				CourseCode: courseCode,
				Name:       courseName,
				Credit:     credit,
				Semester:   semester,
				Status:     1, // 已开课
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			err = s.courseRepo.Create(course)
			if err != nil {
				return err
			}
		}

		// 创建学生选课记录
		studentCourse := &model.StudentCourse{
			StudentID: studentID,
			CourseID:  course.ID,
			Status:    1, // 已选
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.courseRepo.CreateStudentCourse(studentCourse)
		if err != nil {
			return err
		}

		// 解析课程安排信息
		weekday, _ := strconv.Atoi(strings.TrimSpace(row[4]))
		startTime := strings.TrimSpace(row[5])
		endTime := strings.TrimSpace(row[6])
		classroom := ""
		building := ""
		if len(row) > 7 {
			classroom = strings.TrimSpace(row[7])
		}
		if len(row) > 8 {
			building = strings.TrimSpace(row[8])
		}

		// 创建课程安排
		schedule := &model.CourseSchedule{
			CourseID:  course.ID,
			TeacherID: 0, // 需要根据实际情况设置
			Classroom: classroom,
			Building:  building,
			Weekday:   weekday,
			StartWeek: 1,  // 默认值，需要根据实际情况调整
			EndWeek:   16, // 默认值，需要根据实际情况调整
			StartTime: startTime,
			EndTime:   endTime,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.courseRepo.CreateCourseSchedule(schedule)
		if err != nil {
			return err
		}
	}

	return nil
}
