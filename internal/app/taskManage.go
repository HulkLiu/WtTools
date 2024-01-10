package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"strconv"
	"sync"
	"time"
)

var err error

func (a *App) GetTaskList() *utils.Response {

	tasks := a.TaskManage.GetTasks()
	//log.Printf("GetTaskList :%+v", tasks)
	return utils.Success(tasks)

}

func (a *App) AddTask(content string) *utils.Response {
	log.Printf("AddTask content ====> %+v", content)

	t := Task{}
	if err := json.Unmarshal([]byte(content), &t); err != nil {
		a.Log.Errorf("AddTask Unmarshal err: %v", err)
		return utils.Fail(err.Error())
	}
	log.Printf("AddTask ====> %v", t)
	_, err := a.TaskManage.NewTask(t.Name, t.Description)
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	return utils.Success("任务添加成功")
}
func (a *App) DeleteTask(id string) *utils.Response {
	// 删除任务
	i, err := strconv.Atoi(id)
	err = a.TaskManage.DeleteTask(i)
	if err != nil {
		return utils.Fail(fmt.Sprintf("Error deleting task:%v", err))
	}
	return utils.Success("任务删除成功")
}

func (a *App) UpdateTask(content string) *utils.Response {
	log.Printf("UpdateTask content ====> %+v", content)

	t := Task{}
	if err := json.Unmarshal([]byte(content), &t); err != nil {
		a.Log.Errorf("UpdateTask Unmarshal err: %v", err)
		return utils.Fail(err.Error())
	}
	err = a.TaskManage.UpdateTask(t.ID, t.Name, t.Description, t.Completed)
	if err != nil {
		return utils.Fail(fmt.Sprintf("Error UpdateTask task:%v", err))
	}
	return utils.Success("任务修改成功")
}

func (a *App) FilterTasks(content string) *utils.Response {

	tasks := a.TaskManage.FilterTasks(func(t Task) int {
		return t.Completed
	})
	return utils.Success(tasks)

}

// Task 表示一个任务
type Task struct {
	ID          int    `json:"Id" gorm:"column:Id;AUTO_INCREMENT;not null"`
	Name        string `json:"Name" gorm:"column:Name;not null"`
	Description string `json:"Description" gorm:"column:Description;not null"`
	Completed   int    `json:"Completed" gorm:"column:Completed;not null"`
	CreatedAt   string `json:"CreatedAt" gorm:"column:CreatedAt;not null"`
	UpdatedAt   string `json:"UpdatedAt" gorm:"column:UpdatedAt;not null"`
}

// TableName ...
func (t *Task) TableName() string {
	return "task"
}

func NewTaskDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{
		Logger: gormLogger.Discard,
	})
	if err != nil {
		log.Fatalf(" mysql connect failed,err:%v", err)
		return nil

	}

	// 判断表是否存在
	if !db.Migrator().HasTable(&Task{}) {
		// 创建表
		err = db.AutoMigrate(&Task{})
		if err != nil {
			return nil
		}
	}
	return db
}

// TaskList 表示任务列表
type TaskList struct {
	Tasks []Task
	mu    sync.RWMutex // 用于并发控制
}

// TaskManager 表示任务管理器，包含任务列表
type TaskManager struct {
	TaskList TaskList
	TasksDB  *gorm.DB
	DSN      string
}

// NewTask 创建一个新任务
func (tm *TaskManager) NewTask(name, description string) (Task, error) {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	task := Task{
		Name:        name,
		Description: description,
		Completed:   0,
		CreatedAt:   time.Now().Format(config.TimeFormat),
		UpdatedAt:   time.Now().Format(config.TimeFormat),
	}

	tm.TaskList.Tasks = append(tm.TaskList.Tasks, task)
	if err := tm.TasksDB.Create(&task).Error; err != nil {
		return Task{}, err
	}
	return task, nil
}

// CompleteTask 标记任务为完成
func (tm *TaskManager) CompleteTask(id int) error {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	for i, task := range tm.TaskList.Tasks {
		if task.ID == id {
			tm.TaskList.Tasks[i].Completed = 1
			tm.TaskList.Tasks[i].UpdatedAt = time.Now().Format(config.TimeFormat)
			return nil
		}
	}

	return errors.New("task not found")
}

// DeleteTask 删除一个任务
func (tm *TaskManager) DeleteTask(id int) error {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	for i, task := range tm.TaskList.Tasks {
		if task.ID == id {
			tm.TaskList.Tasks = append(tm.TaskList.Tasks[:i], tm.TaskList.Tasks[i+1:]...)
			return nil
		}
	}

	return errors.New("task not found")
}

// UpdateTask 更新一个任务
func (tm *TaskManager) UpdateTask(id int, name, description string, completed int) error {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	for i, task := range tm.TaskList.Tasks {
		if task.ID == id {
			tm.TaskList.Tasks[i].Name = name
			tm.TaskList.Tasks[i].Description = description
			tm.TaskList.Tasks[i].UpdatedAt = time.Now().Format(config.TimeFormat)
			tm.TaskList.Tasks[i].Completed = completed
			return nil
		}
	}

	return errors.New("task not found")
}

// FilterTasks 返回符合过滤条件的任务
func (tm *TaskManager) FilterTasks(filterFunc func(Task) int) []Task {
	tm.TaskList.mu.RLock()
	defer tm.TaskList.mu.RUnlock()

	var filteredTasks []Task
	for _, task := range tm.TaskList.Tasks {
		if filterFunc(task) == 1 {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}

// GetTasks 返回所有任务
func (tm *TaskManager) GetTasks() []Task {
	tm.TaskList.mu.RLock()
	defer tm.TaskList.mu.RUnlock()

	List := make([]Task, 0)

	tm.TasksDB.Order("CreatedAt DESC").Find(&List)
	return List
}

// GetTaskByID 返回指定ID的任务
func (tm *TaskManager) GetTaskByID(id int) (Task, error) {
	tm.TaskList.mu.RLock()
	defer tm.TaskList.mu.RUnlock()

	for _, task := range tm.TaskList.Tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, errors.New("task not found")
}
