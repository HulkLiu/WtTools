package app

import (
	"fmt"
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/parser"
	"github.com/HulkLiu/WtTools/internal/persist"
	"github.com/HulkLiu/WtTools/internal/scheduler"
	"github.com/HulkLiu/WtTools/internal/utils"
	"github.com/go-redis/redis"
	"log"
	"regexp"
	"sync"
	"time"
)

type SearchEngine interface {
	Search(keyword string) error
}

func getMysqlData() []string {
	//mysql 查询得出结果
	arr := []string{
		"https://www.666php.com/?s=",
		"https://www.600xue.com/?s=",
	}
	return arr
}

func (a *App) SearchCourse(keyword string) *utils.Response {
	keyword = "golang"

	//Redis 缓存
	cachedResult, err := a.DB.RedisClient.Get(keyword).Result()
	if err == redis.Nil {

		itemChan, err := persist.ItemSaver(config.ElasticIndex, a.DB.EsClient)
		if err != nil {
			log.Printf("initEsData Connect failed,err%v\n", err)
		}

		//初始化爬虫引擎
		e := engine.Concurrent{
			Scheduler:        &scheduler.QueuedScheduler{},
			WorkerCount:      2,
			ItemChan:         itemChan,
			RequestProcessor: engine.Worker,
		}
		//初始化爬虫对象
		searchEngines := map[string]SearchEngine{
			"666php": &parser.Java666SearchEngine{E: e},
			"600xue": &parser.Xue600SearchEngine{E: e},
		}
		//执行爬虫动作
		go searchManage(keyword, getMysqlData(), searchEngines)
	} else if err != nil {
		log.Printf("Failed to get from cache: %v", err)
	} else {
		//查到到缓存 直接返回缓存数据
		return utils.Success(cachedResult)
	}

	//查询ES
	list, err := a.CourseManage.CourseList(keyword)
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	//设置 Redis缓存
	err = a.DB.RedisClient.Set(keyword, list.Items, 10*time.Second).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}
	log.Printf("SearchCourse ==> %v", len(list.Items))
	return utils.Success(list.Items)
}
func regexpUrl(url string) string {
	rule := `https://www\.(.*?)\.com/`

	// 编译正则表达式
	reg := regexp.MustCompile(rule)

	// 提取匹配的字符串
	result := reg.FindStringSubmatch(url)
	return result[1]
}

var wg sync.WaitGroup

func searchManage(keyword string, arr []string, searchEngines map[string]SearchEngine) {
	if len(arr) == 0 {
		return
	}

	for _, v := range arr {
		wg.Add(1)
		url := v + keyword
		go func(u string) {
			defer wg.Done()
			searchEngine := searchEngines[regexpUrl(url)]
			searchEngine.Search(u)
		}(url)
	}

	wg.Wait()
}
