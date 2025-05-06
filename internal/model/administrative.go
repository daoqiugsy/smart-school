package model

import (
	"time"
)

// LeaveApplication 请假申请
type LeaveApplication struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Type        int       `gorm:"not null" json:"type"`             // 0:事假 1:病假 2:其他
	Reason      string    `gorm:"type:text;not null" json:"reason"` // 请假原因
	StartTime   time.Time `json:"start_time"`                       // 开始时间
	EndTime     time.Time `json:"end_time"`                         // 结束时间
	Attachments string    `gorm:"type:text" json:"attachments"`     // 附件路径，多个用逗号分隔
	Status      int       `gorm:"not null;default:0" json:"status"` // 0:待审核 1:已批准 2:已拒绝
	ApproverID  uint      `gorm:"default:null" json:"approver_id"`  // 审批人ID
	Comment     string    `gorm:"type:text" json:"comment"`         // 审批意见
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	Approver    User      `gorm:"foreignKey:ApproverID" json:"approver"`
}

// TableName 指定表名
func (LeaveApplication) TableName() string {
	return "leave_applications"
}

// ReimbursementApplication 报销申请
type ReimbursementApplication struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Title       string    `gorm:"size:100;not null" json:"title"`        // 报销标题
	Amount      float64   `gorm:"not null" json:"amount"`                // 报销金额
	Description string    `gorm:"type:text;not null" json:"description"` // 报销说明
	Attachments string    `gorm:"type:text;not null" json:"attachments"` // 附件路径，多个用逗号分隔
	Status      int       `gorm:"not null;default:0" json:"status"`      // 0:待审核 1:已批准 2:已拒绝 3:已报销
	ApproverID  uint      `gorm:"default:null" json:"approver_id"`       // 审批人ID
	Comment     string    `gorm:"type:text" json:"comment"`              // 审批意见
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	Approver    User      `gorm:"foreignKey:ApproverID" json:"approver"`
}

// TableName 指定表名
func (ReimbursementApplication) TableName() string {
	return "reimbursement_applications"
}

// AssetApplication 资产申请
type AssetApplication struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	AssetName  string    `gorm:"size:100;not null" json:"asset_name"` // 资产名称
	Quantity   int       `gorm:"not null" json:"quantity"`            // 申请数量
	Purpose    string    `gorm:"type:text;not null" json:"purpose"`   // 用途说明
	StartTime  time.Time `json:"start_time"`                          // 开始使用时间
	EndTime    time.Time `json:"end_time"`                            // 结束使用时间
	Status     int       `gorm:"not null;default:0" json:"status"`    // 0:待审核 1:已批准 2:已拒绝 3:已归还
	ApproverID uint      `gorm:"default:null" json:"approver_id"`     // 审批人ID
	Comment    string    `gorm:"type:text" json:"comment"`            // 审批意见
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Approver   User      `gorm:"foreignKey:ApproverID" json:"approver"`
}

// TableName 指定表名
func (AssetApplication) TableName() string {
	return "asset_applications"
}
