-- 作业信息表
CREATE TABLE `assignments` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `course_id` bigint(20) UNSIGNED NOT NULL,
  `title` varchar(100) NOT NULL COMMENT '作业标题',
  `description` text DEFAULT NULL COMMENT '作业描述',
  `deadline` datetime DEFAULT NULL COMMENT '截止时间',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0:已取消 1:进行中 2:已截止',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  CONSTRAINT `fk_assignments_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作业信息';

-- 学生作业提交记录表
CREATE TABLE `student_assignments` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `student_id` bigint(20) UNSIGNED NOT NULL,
  `assignment_id` bigint(20) UNSIGNED NOT NULL,
  `content` text DEFAULT NULL COMMENT '提交内容',
  `attachments` text DEFAULT NULL COMMENT '附件路径，多个用逗号分隔',
  `score` float DEFAULT NULL COMMENT '分数',
  `comment` text DEFAULT NULL COMMENT '评语',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:未提交 1:已提交 2:已批改',
  `submit_time` datetime DEFAULT NULL COMMENT '提交时间',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_assignment_id` (`assignment_id`),
  CONSTRAINT `fk_student_assignments_student` FOREIGN KEY (`student_id`) REFERENCES `students` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_student_assignments_assignment` FOREIGN KEY (`assignment_id`) REFERENCES `assignments` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学生作业提交记录';

-- 考试信息表
CREATE TABLE `exams` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `course_id` bigint(20) UNSIGNED NOT NULL,
  `title` varchar(100) NOT NULL COMMENT '考试标题',
  `description` text DEFAULT NULL COMMENT '考试描述',
  `location` varchar(100) DEFAULT NULL COMMENT '考试地点',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `exam_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:平时测验 1:期中考试 2:期末考试',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:未开始 1:进行中 2:已结束',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_course_id` (`course_id`),
  CONSTRAINT `fk_exams_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='考试信息';