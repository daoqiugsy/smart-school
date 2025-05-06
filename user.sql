-- 用户基础信息表
CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(100) NOT NULL,
  `real_name` varchar(50) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `user_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:学生 1:教师 2:管理员',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0:禁用 1:启用',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基础信息';

-- 学生信息表
CREATE TABLE `students` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `student_id` varchar(20) NOT NULL COMMENT '学号',
  `grade` varchar(20) DEFAULT NULL COMMENT '年级',
  `class` varchar(50) DEFAULT NULL COMMENT '班级',
  `major` varchar(50) DEFAULT NULL COMMENT '专业',
  `department` varchar(50) DEFAULT NULL COMMENT '院系',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  UNIQUE KEY `idx_student_id` (`student_id`),
  CONSTRAINT `fk_students_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生信息';

-- 教师信息表
CREATE TABLE `teachers` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `teacher_id` varchar(20) NOT NULL COMMENT '工号',
  `title` varchar(50) DEFAULT NULL COMMENT '职称',
  `department` varchar(50) DEFAULT NULL COMMENT '院系',
  `office` varchar(50) DEFAULT NULL COMMENT '办公室',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  UNIQUE KEY `idx_teacher_id` (`teacher_id`),
  CONSTRAINT `fk_teachers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教师信息';