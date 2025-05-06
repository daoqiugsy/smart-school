-- 通知信息表
CREATE TABLE `notifications` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL COMMENT '通知标题',
  `content` text NOT NULL COMMENT '通知内容',
  `type` tinyint(4) NOT NULL COMMENT '0:系统通知 1:课程通知 2:考试通知 3:作业通知 4:行政通知',
  `priority` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0:低 1:中 2:高',
  `source_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '来源ID，如课程ID、考试ID等',
  `source_type` varchar(50) DEFAULT NULL COMMENT '来源类型，如Course、Exam等',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知信息';

-- 用户通知关联表
CREATE TABLE `user_notifications` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `notification_id` bigint(20) UNSIGNED NOT NULL,
  `is_read` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已读',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `delivery_status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:待发送 1:已发送 2:发送失败',
  `delivery_type` tinyint(4) NOT NULL COMMENT '0:APP 1:短信 2:邮箱 3:语音',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_notification_id` (`notification_id`),
  CONSTRAINT `fk_user_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_notifications_notification` FOREIGN KEY (`notification_id`) REFERENCES `notifications` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户通知关联';