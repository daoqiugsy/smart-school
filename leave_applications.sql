-- 请假申请表
CREATE TABLE `leave_applications` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `type` tinyint(4) NOT NULL COMMENT '0:事假 1:病假 2:其他',
  `reason` text NOT NULL COMMENT '请假原因',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `attachments` text DEFAULT NULL COMMENT '附件路径，多个用逗号分隔',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:待审核 1:已批准 2:已拒绝',
  `approver_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '审批人ID',
  `comment` text DEFAULT NULL COMMENT '审批意见',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_approver_id` (`approver_id`),
  CONSTRAINT `fk_leave_applications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_leave_applications_approver` FOREIGN KEY (`approver_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='请假申请';

-- 报销申请表
CREATE TABLE `reimbursement_applications` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `title` varchar(100) NOT NULL COMMENT '报销标题',
  `amount` decimal(10,2) NOT NULL COMMENT '报销金额',
  `description` text NOT NULL COMMENT '报销说明',
  `attachments` text NOT NULL COMMENT '附件路径，多个用逗号分隔',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:待审核 1:已批准 2:已拒绝 3:已报销',
  `approver_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '审批人ID',
  `comment` text DEFAULT NULL COMMENT '审批意见',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_approver_id` (`approver_id`),
  CONSTRAINT `fk_reimbursement_applications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_reimbursement_applications_approver` FOREIGN KEY (`approver_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报销申请';

-- 资产申请表
CREATE TABLE `asset_applications` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `asset_name` varchar(100) NOT NULL COMMENT '资产名称',
  `quantity` int(11) NOT NULL COMMENT '申请数量',
  `purpose` text NOT NULL COMMENT '用途说明',
  `start_time` datetime DEFAULT NULL COMMENT '开始使用时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束使用时间',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:待审核 1:已批准 2:已拒绝 3:已归还',
  `approver_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '审批人ID',
  `comment` text DEFAULT NULL COMMENT '审批意见',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_approver_id` (`approver_id`),
  CONSTRAINT `fk_asset_applications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_asset_applications_approver` FOREIGN KEY (`approver_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资产申请';