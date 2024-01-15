package app

import (
	"context"
	"fmt"
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/parser"
	"github.com/HulkLiu/WtTools/internal/persist"
	"github.com/HulkLiu/WtTools/internal/scheduler"
	"github.com/HulkLiu/WtTools/internal/utils"
	"golang.org/x/sync/errgroup"
	"log"
	"regexp"
)

func (a *App) SearchCourse(keyword string) *utils.Response {
	keyword = "golang"

	searchManage(keyword)

	list, err := a.CourseManage.CourseList(keyword)
	log.Printf("SearchCourse ==> %v", list.Items)
	if err != nil {
		return utils.Fail(fmt.Sprintf("%v", err))
	}
	return utils.Success(list.Items)
}
func regexpUrl(url string) string {
	rule := `https://www\.(.*?)\.com/\?s=`

	// 编译正则表达式
	reg := regexp.MustCompile(rule)

	// 提取匹配的字符串
	result := reg.FindStringSubmatch(url)
	return result[1]
}
func searchManage(keyword string) {
	//mysql 查询得出结果
	arr := []string{
		//"https://www.666php.com/?s=",
		"https://www.600xue.com/?s=",
	}

	itemChan, err := persist.ItemSaver(config.ElasticIndex)
	if err != nil {
		log.Printf("initEsData Connect failed,err%v\n", err)
	}

	e := engine.Concurrent{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      1,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(1)

	for _, v := range arr {
		url := v + keyword

		g.Go(func() error {
			err = doSearch(url, e)
			if err != nil {
				return err
			}
			return nil
		})
	}
	_ = g.Wait()

}

func doSearch(url string, e engine.Concurrent) error {
	switch regexpUrl(url) {
	case "666php":
		Java666 := parser.Java666{
			ClassListRe: regexp.MustCompile(`<a target="_blank" href="(https:\/\/www.666php.com\/\d{1,}.html)" title=".*?" rel="bookmark">(.*?)<\/a>`),

			ClassPageRe: regexp.MustCompile(`(https:\/\/www.666php.com)`),
			UrlMap:      make(map[string]string),

			TitleRE: regexp.MustCompile(`<h2>(.*)<\/h2>`),
			IdRe:    regexp.MustCompile(`com\/(.*?).html`),
		}
		e.Run(engine.Request{
			Url:        url,
			ParserFunc: Java666.ClassList,
		})
	case "600xue":
		xue600 := parser.Xue600{
			ClassListRe: regexp.MustCompile(`<a target="_blank" href="(.*?)" title=".*?" rel="bookmark">(.*?)<\/a>`),

			ClassPageRe: regexp.MustCompile(`(https:\/\/www.600xue.com)`),
			UrlMap:      make(map[string]string),

			TitleRE: regexp.MustCompile(`<h2>(.*)<\/h2>`),
			IdRe:    regexp.MustCompile(`com\/(.*?).html`),
		}
		e.Run(engine.Request{
			Url:        url,
			ParserFunc: xue600.ClassList,
		})
	}
	return nil
}
