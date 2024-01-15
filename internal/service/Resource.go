package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Resource 表示一个资源
type Resource struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	URL       string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ResourceRule  表示一个规则
type ResourceRule struct {
	ID          uint `gorm:"primaryKey"`
	Description string
}

// ResourceLog  表示一个日志条目
type ResourceLog struct {
	ID        uint `gorm:"primaryKey"`
	Message   string
	Timestamp time.Time
}

// ResourcesList  资源列表
type ResourcesList struct {
	List *[]Resource
	//List *  map[string]Resource
	mu sync.RWMutex // 用于并发控制
}
type ResourceManage struct {
	Resources    ResourcesList
	RuleResource ResourceRule
	LogResource  ResourceLog
	db           *gorm.DB
	DSN          string
}

func NewResource(db *gorm.DB) ResourceManage {

	r := ResourceManage{
		db: db,
	}
	// 设置创建表的默认字符集为utf-8
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8")
	if err := db.AutoMigrate(&Resource{}, &ResourceRule{}, &ResourceLog{}); err != nil {
		log.Fatal("NewResourceDB:Failed to AutoMigrate to database:", err)
	}

	_ = r.ListResources()

	return r
}

// ListResources 列出所有资源
func (r *ResourceManage) ListResources() error {
	var resources []Resource
	//var resources map[string]Resource
	result := r.db.Find(&resources)
	r.Resources.List = &resources
	return result.Error
}

// ImportResourcesFromExcel 从Excel文件导入资源到数据库
func (r *ResourceManage) ImportResourcesFromExcel(filename string) error {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if row.Cells != nil {
				var name, url, status string
				if len(row.Cells) >= 1 {
					name = row.Cells[0].String()
				}
				if len(row.Cells) >= 2 {
					url = row.Cells[1].String()
				}
				if len(row.Cells) >= 3 {
					status = row.Cells[2].String()
				}
				data := Resource{
					Name:   name,
					URL:    url,
					Status: status,
				}
				r.AddResource(data)
			}
		}
	}
	return nil
}

// ExportResourcesToExcel 导出数据库中的资源到Excel文件
func (r *ResourceManage) ExportResourcesToExcel(filename string) error {
	resources := r.Resources.List

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err := file.AddSheet("Resources")
	if err != nil {
		return err
	}

	for _, resource := range *resources {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = resource.Name
		cell = row.AddCell()
		cell.Value = resource.URL
		cell = row.AddCell()
		cell.Value = resource.Status
	}

	err = file.Save(filename)
	if err != nil {
		return err
	}
	fmt.Printf("Exported data to %s\n", filename)
	return nil
}

// AddResource 创建新资源并插入数据库
func (r *ResourceManage) AddResource(t Resource) error {
	resource := Resource{
		Name:      strings.Trim(t.Name, " "),
		URL:       strings.Trim(t.URL, " "),
		Status:    strings.Trim(t.Status, " "),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, v := range *r.Resources.List {
		if v.URL == resource.URL {
			return errors.New("URL 已存在")
		}
	}

	result := r.db.Create(&resource)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateResource 更新资源信息
func (r *ResourceManage) UpdateResource(t Resource) error {

	result := r.db.Model(&Resource{}).Where("id = ?", t.ID).Updates(Resource{URL: t.URL, Status: t.Status, UpdatedAt: time.Now()})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteResource 删除资源
func (r *ResourceManage) DeleteResource(id int) error {
	result := r.db.Delete(&Resource{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CheckResourceAll 检查资源链接是否可用
func (r *ResourceManage) CheckResourceAll() (map[string]string, error) {
	var resources []Resource
	//result := db.First(&resource, id)
	result := r.db.Find(&resources)
	var checkMap = make(map[string]string, len(resources))

	if result.Error != nil {
		log.Println("Error finding service:", result.Error)
		return checkMap, result.Error
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(5)

	for _, v := range resources {
		resource := v

		g.Go(func() error {
			client := http.Client{
				Timeout: 3 * time.Second,
			}
			_, err := client.Get(resource.URL)
			if err != nil {
				checkMap[resource.URL] = fmt.Sprintf("%v", err)
			}
			checkMap[resource.URL] = "true"
			return nil
		})
	}
	err := g.Wait()
	return checkMap, err
	//defer resp.Body.Close()

}

func test() {
	r := ResourceManage{}
	// 示例：创建资源并添加到数据库

	// 示例：更新资源
	//r.UpdateResource(1, "https://www.example.org", "Updated")

	// 示例：删除资源
	r.DeleteResource(1)

	// 示例：检查资源
	//isAvailable := r.CheckResource(1)
	//fmt.Println("Resource available:", isAvailable)

	// 示例：添加日志
	//AddLog("Resource ExampleResource was checked for availability")

	// 示例：列出所有资源
	err := r.ListResources()
	if err != nil {
		log.Fatal("Error listing resources:", err)
	}

	// 示例：从Excel导入资源
	err = r.ImportResourcesFromExcel("resources.xlsx")
	if err != nil {
		log.Fatal("Error importing resources from Excel:", err)
	}

	// 示例：导出资源到Excel
	err = r.ExportResourcesToExcel("exported_resources.xlsx")
	if err != nil {
		log.Fatal("Error exporting resources to Excel:", err)
	}
}
