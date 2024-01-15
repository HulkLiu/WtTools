package app

import (
	"encoding/json"
	"fmt"
	"github.com/HulkLiu/WtTools/internal/service"
	"github.com/HulkLiu/WtTools/internal/utils"
	"log"
	"strconv"
)

var err error

func (a *App) GetTaskList() *utils.Response {

	tasks := a.TaskManage.GetTasks()
	//log.Printf("GetTaskList :%+v", tasks)
	return utils.Success(tasks)

}

func (a *App) AddTask(content string) *utils.Response {
	log.Printf("AddTask content ====> %+v", content)

	t := service.Task{}
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

	t := service.Task{}
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

	tasks := a.TaskManage.FilterTasks(func(t service.Task) int {
		return t.Completed
	})
	return utils.Success(tasks)

}
