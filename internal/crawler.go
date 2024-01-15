package main

import (
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/parser"
	"github.com/HulkLiu/WtTools/internal/persist"
	"github.com/HulkLiu/WtTools/internal/scheduler"
	"github.com/olivere/elastic/v7"
	"log"
)

func main() {
	url := "https://666php.com/"
	//url := "https://shipin.vjshi.com/piantou/2?direct=gJLngJeX0czIwICXl4CW5JSTnJyUkZeUkJCSlpCAl+aAl5fM1uHM18DG0YCXl4CW5NHX0MCAkuE"
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Printf("err:%v", err)
	}
	itemChan, err := persist.ItemSaver(config.ElasticIndex, client)
	if err != nil {
		log.Printf("initEsData Connect failed,err%v\n", err)
		return
	}

	e := engine.Concurrent{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      5,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}
	a := parser.Java666{}
	e.Run(engine.Request{
		Url:        url,
		ParserFunc: a.ClassList,
	})

}
