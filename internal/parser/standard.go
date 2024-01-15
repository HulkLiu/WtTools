package parser

import (
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/fetcher"
	"github.com/HulkLiu/WtTools/internal/model"
	"log"
	"regexp"
	"time"
)

func NewObjSearchEngine(e engine.Concurrent) *ObjSearchEngine {
	return &ObjSearchEngine{
		E:      e,
		parser: NewObj(),
	}
}

type ObjSearchEngine struct {
	E      engine.Concurrent
	parser *Obj
}
type Obj struct {
	ClassListRe *regexp.Regexp

	ClassPageRe *regexp.Regexp
	UrlMap      map[string]string

	TitleRE *regexp.Regexp
	IdRe    *regexp.Regexp

	LanguageTypeSel string
	CatalogSel      string
	LastTimeSel     string
	LimitItem       int
}

func NewObj() *Obj {
	return &Obj{
		ClassListRe: regexp.MustCompile(`<a target="_blank" href="(https:\/\/www.666php.com\/\d{1,}.html)" title=".*?" rel="bookmark">(.*?)<\/a>`),
		ClassPageRe: regexp.MustCompile(`(https:\/\/www.666php.com)`),
		UrlMap:      make(map[string]string),
		TitleRE:     regexp.MustCompile(`<h2>(.*)<\/h2>`),
		IdRe:        regexp.MustCompile(`com\/(.*?).html`),
		LimitItem:   5,
	}
}

func (j *ObjSearchEngine) Search(url string) error {
	log.Printf("ObjSearchEngine: Searching for %s", url)
	j.E.Run(engine.Request{
		Url:        url,
		ParserFunc: j.parser.ClassPage,
	})
	return nil
}

func (a *Obj) ClassPage(c []byte) engine.ParseResult {
	result := engine.ParseResult{}
	matches := a.ClassListRe.FindAllSubmatch(c, -1)

	for k, m := range matches {
		if a.LimitItem > 0 && k == a.LimitItem {
			break
		}
		if len(m) <= 2 {
			continue
		}
		title := string(m[2])
		url := string(m[1])

		if a.isUrlExists(url) {
			log.Println(url)
			continue
		}

		log.Printf(" ClassPageB url :%v,title:%v", url, title)

		result.Requests = append(result.Requests, engine.Request{
			Url: url,

			ParserFunc: func(c []byte) engine.ParseResult {
				return a.ParseProfile(c, url, title)
			},
		})
	}

	return result

}

func (a *Obj) ClassList(c []byte) engine.ParseResult {
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

		if a.isUrlExists(url) {
			log.Println(url)
			continue
		}
		//fmt.Printf("Find %v ClassPageA url %v\n", nu, url)
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

// TrimHtml 去除字符串中的html标签

func (a *Obj) ParseProfile(c []byte, url string, title string) engine.ParseResult {
	p := model.Profile{}
	p.Title = title
	if title == "" {
		p.Title = a.extractString(c, a.TitleRE)
	}

	//分类
	p.LanguageType = fetcher.GetTextByQuery(string(c), "body > div > div.site-content > div > div.breadcrumbs > a:nth-child(5)")

	//描述
	p.Catalog = fetcher.GetTextByQuery(string(c), "#post-"+a.extractString([]byte(url), a.IdRe)+" > div.container > div.entry-wrapper > article > div > pre")

	//时间
	p.LastTime = fetcher.GetTextByQuery(string(c), "body > div > div.site-content > div > section > div > hgroup > div.meta > div.description > span")

	p.CreateAt = time.Now().Format(config.DateFormat)
	p.Status = "待确认"
	//log.Printf("%+v", p)
	//utils.DownFile(fmt.Sprintf("%+v", p))

	res := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "php666",
				Id:      a.extractString([]byte(url), a.IdRe),
				Payload: p,
			},
		},
	}

	//fmt.Printf("%+v\n", res)
	return res
}

func (a *Obj) extractString(c []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(c)

	//for k, v := range match {
	//	fmt.Printf("%v -> %s\n", k, v)
	//}

	if len(match) > 1 {
		return string(match[1])
	} else {
		return ""
	}
}
func (a *Obj) isUrlExists(url string) bool {
	if _, ok := a.UrlMap[url]; ok {
		return true
	}
	a.UrlMap[url] = url
	return false
}
