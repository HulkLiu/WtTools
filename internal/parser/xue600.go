package parser

import (
	"fmt"
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/fetcher"
	"github.com/HulkLiu/WtTools/internal/model"
	"log"
	"regexp"
	"time"
)

type Xue600 struct {
	ClassListRe *regexp.Regexp
	ClassPageRe *regexp.Regexp
	UrlMap      map[string]string
	TitleRE     *regexp.Regexp
	IdRe        *regexp.Regexp
}

func (a *Xue600) ClassPage(c []byte) engine.ParseResult {
	result := engine.ParseResult{}
	matches := a.ClassListRe.FindAllSubmatch(c, -1)

	for k, m := range matches {
		if true {
			if k == 3 {
				break
			}
		}
		if len(m) <= 2 {
			continue
		}
		title := string(m[2])
		url := string(m[1])

		if _, ok := a.UrlMap[url]; ok {
			log.Println(url)
			continue
		}
		a.UrlMap[url] = url

		log.Printf("ClassPage => url :%v,title:%v", url, title)

		result.Requests = append(result.Requests, engine.Request{
			Url: url,

			ParserFunc: func(c []byte) engine.ParseResult {
				return a.ParseProfile(c, url, title)
			},
		})
	}

	return result

}

func (a *Xue600) ClassList(c []byte) engine.ParseResult {
	result := engine.ParseResult{}
	matches := a.ClassPageRe.FindAllSubmatch(c, -1)
	//log.Printf("%s", matches)
	//return result
	for k, m := range matches {
		if k == 1 {
			break
		}
		//title := string(m[2])
		url := string(m[1])
		if _, ok := a.UrlMap[url]; ok {
			continue
		}
		a.UrlMap[url] = url

		fmt.Printf("ClassList => ClassPageA url %v\n", url)
		result.Requests = append(result.Requests, engine.Request{
			Url: url,
			ParserFunc: func(c []byte) engine.ParseResult {
				return a.ClassPage(c)
			},
		})
	}

	//println(len(matches))
	return result
}

func (a *Xue600) ParseProfile(c []byte, url string, title string) engine.ParseResult {
	p := model.Profile{}
	p.Title = title
	if title == "" {
		p.Title = extractString(c, a.TitleRE)
	}

	//分类
	p.LanguageType = fetcher.GetTextByQuery(string(c), "body > div > div.site-content > div > div.breadcrumbs > a:nth-child(5)")

	//描述
	p.Catalog = fetcher.GetTextByQuery(string(c), "#post-"+extractString([]byte(url), a.IdRe)+" > div.container > div.entry-wrapper > article > div > pre")

	//时间
	p.LastTime = fetcher.GetTextByQuery(string(c), "body > div > div.site-content > div > section > div > hgroup > div.meta > div.description > span")

	p.CreateAt = time.Now().Format(config.DateFormat)
	p.Status = "待确认"
	log.Printf("%+v", p)
	//utils.DownFile(fmt.Sprintf("%+v", p))

	res := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "php666",
				Id:      extractString([]byte(url), a.IdRe),
				Payload: p,
			},
		},
	}

	//fmt.Printf("%+v\n", res)
	return res
}
