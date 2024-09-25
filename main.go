package main

import (
	"calendar-remind-service/db"
	"calendar-remind-service/task"
	"calendar-remind-service/validate"
	_ "calendar-remind-service/validate"
	"calendar-remind-service/ws"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")
	ws.AddClient(conn) // 添加客户端
	for {
		_, _, err := conn.ReadMessage() // 持续保持连接
		if err != nil {
			ws.RemoveClient(conn) // 连接关闭时移除
			break
		}
	}

}

// 全局变量
var taskManager task.TaskManager

func main() {
	// 加载db.env文件
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("load the .env file fail!: %v", err)
	}
	//初始化连接池
	initDB := db.InitDB()
	//初始化任务管理器
	taskManager = *task.NewTaskManager()
	// 加载任务
	if err := taskManager.LoadTasks(initDB); err != nil {
		panic("failed to load tasks: " + err.Error())
	}
	r := gin.Default()
	// 拦截器
	r.Use(validate.TokenMiddleware())
	// Restful风格
	r.POST("/reminders", CreateReminder)
	r.GET("/reminders", GetUserReminders)
	r.DELETE("/reminders/:id", DeleteReminder)
	r.PUT("/reminders/:id", UpdateUserReminder)
	// 查看当前任务列表id
	r.GET("task", taskManager.GetCurrentTasks)
	// WebSocket 路由
	r.GET("/ws", handleWebSocket)
	err := r.Run(":8080")
	if err != nil {
		return
	}

}

func CreateReminder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	// 转成go结构体
	var reminder db.UserReminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if reminder.SendType == 0 || reminder.ContactInfo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "接收方式和联系方式都不能为空!"})
		return
	}

	reminder.CreatorID = userID.(int)
	db.AddUserReminder(&reminder)
	// 新增一个任务
	addTask := task.Task{
		ID:         reminder.ID,
		ReminderAt: reminder.ReminderAt,
	}
	taskManager.AddTask(addTask)
	c.JSON(http.StatusOK, reminder)
}

func GetUserReminders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	reminders, err := db.GetUserReminderByCreateID(userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reminders)
}

func DeleteReminder(c *gin.Context) {
	userID, exists := c.Get("userID")
	id := c.Param("id")
	// 将 ID 从字符串转换为整数
	ID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	reminder, _ := db.GetUserReminderByID(ID)
	if reminder.CreatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete the remind that is not belong to you"})
		return
	}
	_, err = db.DeleteUserReminder(ID)
	if err != nil {
		return
	}
	taskManager.DeleteTask(reminder.ID)
	c.JSON(http.StatusOK, gin.H{"success": "delete successful"})
}

func UpdateUserReminder(c *gin.Context) {
	userID, exists := c.Get("userID")
	id := c.Param("id")
	// 将 ID 从字符串转换为整数
	ID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	reminder, _ := db.GetUserReminderByID(ID)
	if reminder.CreatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot update the remind that is not belong to you"})
		return
	}
	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新数据库
	if err := db.UpdateUserReminder(ID, reminder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reminder"})
		return
	}
	updatedTask := task.Task{
		ID:         reminder.ID,
		ReminderAt: reminder.ReminderAt,
	}
	taskManager.UpdateTask(updatedTask)
	c.JSON(http.StatusOK, reminder) // 返回更新后的记录
}
