package parser

import (
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/persist"
	"github.com/HulkLiu/WtTools/internal/scheduler"
	"github.com/olivere/elastic/v7"
	"log"
	"regexp"
	"testing"
)

func TestXue600(t *testing.T) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Printf("err:%v", err)
	}
	itemChan, err := persist.ItemSaver(config.ElasticIndex, client)
	if err != nil {
		t.Errorf("initEsData Connect failed,err%v\n", err)
	}

	e := engine.Concurrent{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      5,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}
	data := Xue600{
		ClassListRe: regexp.MustCompile(`<a target="_blank" href="(.*?)" title=".*?" rel="bookmark">(.*?)<\/a>`),

		ClassPageRe: regexp.MustCompile(`(https:\/\/www.600xue.com)`),
		UrlMap:      make(map[string]string),

		TitleRE:         regexp.MustCompile(`<h2>(.*)<\/h2>`),
		IdRe:            regexp.MustCompile(`com\/(.*?)\/`),
		LanguageTypeSel: "body > div > div.site-content > div > div.breadcrumbs > a:nth-child(3)",
	}
	e.Run(engine.Request{
		Url: "https://www.600xue.com/?s=golang",

		ParserFunc: data.ClassList,
	})
}
func extractString(c []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(c)

	if len(match) > 1 {
		return string(match[1])
	} else {
		return ""
	}
}
