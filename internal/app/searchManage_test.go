package app

import (
	"log"
	"testing"
	"time"

	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/persist"
	"github.com/HulkLiu/WtTools/internal/scheduler"
)

func TestSearchManage(t *testing.T) {
	itemChan, err := persist.ItemSaver(config.ElasticIndex, nil)
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
	// searchEngines := map[string]SearchEngine{
	// 	"666php": &parser.Java666SearchEngine{E: e},
	// 	"600xue": &parser.Xue600SearchEngine{E: e},
	// }
	//执行爬虫动作
	go searchManage(e, "", getMysqlData())

	time.Sleep(time.Hour)

}
