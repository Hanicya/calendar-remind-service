package db

import (
	"time"
)

// AddUserReminder 新增
func AddUserReminder(reminder *UserReminder) int {
	result := db.Create(reminder)
	return int(result.RowsAffected)
}

// GetUserReminderByCreateID 获取用户提醒
func GetUserReminderByCreateID(createId int) ([]UserReminder, error) {
	var reminders []UserReminder
	result := db.Where("creator_id= ?", createId).Where("deleted", 0).Find(&reminders)
	return reminders, result.Error
}

func GetUserReminderByID(ID int) (UserReminder, error) {
	var reminder UserReminder
	result := db.Where("id= ?", ID).Find(&reminder)
	return reminder, result.Error
}

// DeleteUserReminder 删除
func DeleteUserReminder(id int) (int, error) {
	if err := db.Model(&UserReminder{}).Where("id = ?", id).Update("deleted", 1).Error; err != nil {
		return -1, nil // 更新失败
	}
	return 0, nil
}

// UpdateUserReminder 根据 ID 更新用户提醒
func UpdateUserReminder(id int, updatedData UserReminder) error {
	updateInfo := UserReminder{
		Content:    updatedData.Content,
		ReminderAt: updatedData.ReminderAt,
	}
	result := db.Model(&UserReminder{}).Where("id = ?", id).Updates(updateInfo)
	return result.Error
}

// UserReminder 用户提醒表映射
type UserReminder struct {
	ID          int `gorm:"primaryKey"` // 主键
	CreatorID   int
	Content     string
	ReminderAt  time.Time
	SendType    int
	ContactInfo string
	Deleted     int `gorm:"default:0" json:"-"` // 逻辑删除字段
}

func (UserReminder) TableName() string {
	return "user_reminder" // 指定表名为 users
}
