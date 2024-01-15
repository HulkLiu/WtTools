package parser

import (
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/persist"
	"github.com/HulkLiu/WtTools/internal/scheduler"
	"regexp"
	"testing"
)

func TestXue600(t *testing.T) {
	itemChan, err := persist.ItemSaver(config.ElasticIndex)
	if err != nil {
		t.Errorf("initEsData Connect failed,err%v\n", err)
	}

	e := engine.Concurrent{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      1,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}

	e.Run(engine.Request{
		Url: "https://www.600xue.com/?s=golang",
		ParserFunc: Xue600{
			ClassListRe: regexp.MustCompile(`<a target="_blank" href="(.*?)" title=".*?" rel="bookmark">(.*?)<\/a>`),

			ClassPageRe: regexp.MustCompile(`(https:\/\/www.600xue.com)`),
			UrlMap:      make(map[string]string),

			TitleRE: regexp.MustCompile(`<h2>(.*)<\/h2>`),
			IdRe:    regexp.MustCompile(`com\/(.*?).html`),
		}.ClassList,
	})
}
