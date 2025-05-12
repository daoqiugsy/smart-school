package main

import (
	"fmt"
	"log"
	"smart-school/internal/handler"
	"smart-school/internal/model"
	"smart-school/internal/repository"
	"smart-school/internal/service"
	"smart-school/pkg/config"
	"smart-school/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置文件
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化JWT
	utils.InitJWT(&cfg.JWT)

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据库表 - 按照依赖关系顺序迁移
	log.Println("开始自动迁移数据库表...")
	// 先迁移没有外键依赖的基础表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("用户表迁移失败: %v", err)
	}

	// 再迁移依赖User表的表
	err = db.AutoMigrate(&model.Student{}, &model.Teacher{})
	if err != nil {
		log.Fatalf("学生/教师表迁移失败: %v", err)
	}

	// 再迁移课程相关表
	err = db.AutoMigrate(&model.Course{})
	if err != nil {
		log.Fatalf("课程表迁移失败: %v", err)
	}

	// 最后迁移依赖多个表的关系表
	err = db.AutoMigrate(&model.CourseSchedule{}, &model.StudentCourse{})
	if err != nil {
		log.Fatalf("课程安排/学生选课表迁移失败: %v", err)
	}

	log.Println("数据库迁移完成")

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	courseRepo := repository.NewCourseRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, studentRepo, teacherRepo)
	scheduleService := service.NewScheduleService(studentRepo, courseRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	// 初始化AI处理器
	aiHandler := handler.NewAIHandler(&cfg.AI.Coze)
	handler.RegisterRoutes(r, authHandler, scheduleHandler, aiHandler)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
