package app

import (
	"encoding/json"
	"fmt"
	"github.com/HulkLiu/WtTools/internal/service"
	"github.com/HulkLiu/WtTools/internal/utils"
	"log"
	"strconv"
)

func (a *App) GetSourceList() *utils.Response {

	err = a.ResourceManage.ListResources()
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	list := a.ResourceManage.Resources.List
	log.Printf("GetSourceList  ====> %+v", list)
	return utils.Success(list)

}
func (a *App) CheckSourceList() *utils.Response {

	data, err := a.ResourceManage.CheckResourceAll()
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	log.Printf("CheckSourceList  ====> %+v", data)
	return utils.Success(data)

}
func (a *App) AddSource(form string) *utils.Response {

	t := service.Resource{}
	if err := json.Unmarshal([]byte(form), &t); err != nil {
		a.Log.Errorf("AddTask Unmarshal err: %v", err)
		return utils.Fail(err.Error())
	}
	log.Printf("AddSource  ====> %+v", t)

	err = a.ResourceManage.AddResource(t)
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	return utils.Success("添加成功")

}
func (a *App) EditSource(form string) *utils.Response {
	t := service.Resource{}
	if err := json.Unmarshal([]byte(form), &t); err != nil {
		a.Log.Errorf("AddTask Unmarshal err: %v", err)
		return utils.Fail(err.Error())
	}
	log.Printf("EditSource  ====> %+v", t)

	err = a.ResourceManage.UpdateResource(t)
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	return utils.Success("编辑成功")

}
func (a *App) DelSource(id string) *utils.Response {
	newId, _ := strconv.Atoi(id)
	log.Printf("DelSource  ====> %+v", newId)

	err = a.ResourceManage.DeleteResource(newId)
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	return utils.Success("删除成功")

}
