-- 课程信息表
CREATE TABLE `courses` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `course_code` varchar(20) NOT NULL COMMENT '课程代码',
  `name` varchar(100) NOT NULL COMMENT '课程名称',
  `description` text DEFAULT NULL COMMENT '课程描述',
  `credit` float NOT NULL DEFAULT 0 COMMENT '学分',
  `semester` varchar(20) NOT NULL COMMENT '学期',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0:未开课 1:已开课 2:已结课',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_course_code` (`course_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程信息';

-- 课程安排表
CREATE TABLE `course_schedules` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `course_id` bigint(20) UNSIGNED NOT NULL,
  `teacher_id` bigint(20) UNSIGNED NOT NULL,
  `classroom` varchar(50) DEFAULT NULL COMMENT '教室',
  `building` varchar(50) DEFAULT NULL COMMENT '教学楼',
  `weekday` tinyint(4) NOT NULL COMMENT '星期几(1-7)',
  `start_week` int(11) NOT NULL COMMENT '开始周次',
  `end_week` int(11) NOT NULL COMMENT '结束周次',
  `start_time` varchar(10) NOT NULL COMMENT '开始时间(HH:MM)',
  `end_time` varchar(10) NOT NULL COMMENT '结束时间(HH:MM)',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  KEY `idx_teacher_id` (`teacher_id`),
  CONSTRAINT `fk_course_schedules_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_course_schedules_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `teachers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程安排';

-- 学生选课记录表
CREATE TABLE `student_courses` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `student_id` bigint(20) UNSIGNED NOT NULL,
  `course_id` bigint(20) UNSIGNED NOT NULL,
  `score` float DEFAULT NULL COMMENT '成绩',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0:退选 1:已选 2:已修完',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_course_id` (`course_id`),
  CONSTRAINT `fk_student_courses_student` FOREIGN KEY (`student_id`) REFERENCES `students` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_student_courses_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生选课记录';