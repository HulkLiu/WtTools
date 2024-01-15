package main

import (
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/parser"
	"github.com/HulkLiu/WtTools/internal/persist"
	"github.com/HulkLiu/WtTools/internal/scheduler"
	"log"
)

func main() {
	url := "https://666php.com/"
	//url := "https://shipin.vjshi.com/piantou/2?direct=gJLngJeX0czIwICXl4CW5JSTnJyUkZeUkJCSlpCAl+aAl5fM1uHM18DG0YCXl4CW5NHX0MCAkuE"

	itemChan, err := persist.ItemSaver(config.ElasticIndex)
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
