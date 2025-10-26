package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"smart-school/internal/handler"
	"smart-school/internal/model"
	"smart-school/internal/repository"
	"smart-school/internal/service"
	"smart-school/pkg/config"
	"smart-school/pkg/logger"
	"smart-school/pkg/utils"
)

func main() {
	// 加载配置文件
	cfg, err := config.Load("config/config_local.yaml")
	//cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	// 2. 初始化日志 (必须是第二步，紧跟在配置加载之后
	if err := logger.InitLogger(cfg.Server.Mode); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Log.Sync() // 确保日志被写入

	// 从这里开始，我们就可以安全地使用 logger.Log 了
	logger.Log.Info("日志系统初始化成功！")
	log.Printf("-------------------------------------------")
	log.Printf("Loaded Coze URL: %s", cfg.AI.Coze.URL)
	log.Printf("Loaded Coze Token: %s", cfg.AI.Coze.Token)
	log.Printf("Loaded Coze WorkflowID: %s", cfg.AI.Coze.WorkflowID)
	log.Printf("-------------------------------------------")
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化JWT
	utils.InitJWT(&cfg.JWT)

	// 初始化数据库连接
	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		cfg.Database.Host = dbHost // 如果环境变量值存在，用环境变量的值覆盖配置文件中的值
	}
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

	// 初始化Redis连接
	redisHost := os.Getenv("REDIS_HOST") // 尝试读取 REDIS_HOST 环境变量
	if redisHost != "" {
		cfg.Redis.Addr = redisHost + ":6379" // 如果环境变量不存在，默认使用 localhost
		logger.Log.Info("检测到 REDIS_HOST 环境变量, Redis 主机设置为", zap.String("host", redisHost))
	}

	rdb := redis.NewClient(&redis.Options{
		//Addr:     redisHost + "6379",
		//Password: "",
		//DB:       0,
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	// 检查Redis连接
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis连接失败: %v", err)
	}
	logger.Log.Info("Redis 连接成功！")

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	courseRepo := repository.NewCourseRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, studentRepo, teacherRepo)
	scheduleService := service.NewScheduleService(studentRepo, courseRepo, rdb)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	// 初始化AI处理器
	aiService := service.NewAIService(&cfg.AI.Coze)
	aiHandler := handler.NewAIHandler(aiService)
	handler.RegisterRoutes(r, authHandler, scheduleHandler, aiHandler)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
