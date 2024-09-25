package task

import (
	"calendar-remind-service/db"
	"calendar-remind-service/notice"
	"calendar-remind-service/ws"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	ID         int
	ReminderAt time.Time
}

// 任务管理器
type TaskManager struct {
	tasks map[int]chan struct{}
	mu    sync.Mutex
}

// 新建任务管理器
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[int]chan struct{}),
	}
}

// 加载任务
func (tm *TaskManager) LoadTasks(db *gorm.DB) error {
	var tasks []Task
	db.Raw("SELECT id,reminder_at FROM user_reminder WHERE deleted = 0 and reminder_at > NOW()").Scan(&tasks)
	// 查询结果集
	fmt.Printf("需要创建的任务数量:%s\n", strconv.Itoa(len(tasks)))
	for _, task := range tasks {
		//tm.tasks[task.ID] = task
		//go tm.scheduleTask(task) // 启动定时任务
		tm.mu.Lock()                            // 确保并发安全
		tm.tasks[task.ID] = make(chan struct{}) // 初始化取消 channel
		tm.mu.Unlock()
		go tm.scheduleTask(task) // 启动定时任务
	}
	return nil
}

// 调度任务
func (tm *TaskManager) scheduleTask(task Task) {
	tm.mu.Lock()
	cancelChan := make(chan struct{}) // 用于取消任务的 channel
	tm.tasks[task.ID] = cancelChan
	tm.mu.Unlock()

	duration := time.Until(task.ReminderAt)
	if duration > 0 {
		fmt.Printf("创建定时任务%d\n", task.ID)
		select {
		case <-time.After(duration): // 任务到期
			fmt.Printf("Executing task id: %d\n", task.ID)
			Reminder, err := db.GetUserReminderByID(task.ID)
			if err != nil {
				log.Printf("the task %d is missing\n", task.ID)
				return
			}
			// 手机发送
			var notifier notice.Notice
			if Reminder.SendType == 1 {
				notifier = &notice.PhoneNotifier{}
				notifier.Notice(Reminder.ContactInfo, Reminder.Content)
				// 广播形式
				ws.SendMessage(Reminder.Content)
				// 邮箱发送
			} else if Reminder.SendType == 2 {
				notifier = &notice.EmailNotifier{}
				notifier.Notice(Reminder.ContactInfo, Reminder.Content)
				// 广播形式
				ws.SendMessage(Reminder.Content)
			} else {
				log.Printf("the task %d send_type is null\n", task.ID)
			}
		case <-cancelChan: // 任务被取消
			fmt.Printf("任务%d已取消\n", task.ID)
		}
	}
	tm.mu.Lock()
	delete(tm.tasks, task.ID) // 确保在任务完成后删除任务
	tm.mu.Unlock()
}

// 更新任务
func (tm *TaskManager) UpdateTask(task Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// 取消当前任务
	if cancelChan, exists := tm.tasks[task.ID]; exists {
		close(cancelChan)         // 关闭取消 channel
		delete(tm.tasks, task.ID) // 删除旧任务
	}
	// 更新任务并重新调度
	tm.tasks[task.ID] = make(chan struct{}) // 创建新的取消 channel
	go tm.scheduleTask(task)                // 重新调度任务
}

// 新增任务
func (tm *TaskManager) AddTask(task Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// 检查任务是否已存在
	if _, exists := tm.tasks[task.ID]; exists {
		fmt.Printf("任务 %d 已存在，无法添加。\n", task.ID)
		return // 任务已存在，返回
	}
	// 添加新任务并启动调度
	tm.tasks[task.ID] = make(chan struct{}) // 创建取消 channel
	go tm.scheduleTask(task)                // 启动调度
}

// 删除任务
func (tm *TaskManager) DeleteTask(taskID int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if cancelChan, exists := tm.tasks[taskID]; exists {
		close(cancelChan)        // 关闭 channel，通知 goroutine 取消任务
		delete(tm.tasks, taskID) // 从任务列表中删除任务
		fmt.Printf("任务%d已取消\n", taskID)
	} else {
		fmt.Printf("任务%d不存在\n", taskID)
	}

}

// 查看当前的任务列表
func (tm *TaskManager) GetCurrentTasks(c *gin.Context) {
	tm.mu.Lock()         // 加锁
	defer tm.mu.Unlock() // 解锁
	taskIDs := make([]int, 0, len(tm.tasks))
	for taskID := range tm.tasks {
		taskIDs = append(taskIDs, taskID) // 添加任务 ID 到切片
	}
	c.JSON(http.StatusOK, taskIDs) // 返回任务 ID 列表
}
