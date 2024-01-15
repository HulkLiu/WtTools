package parser

import (
	"log"
	"regexp"
	"time"

	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/fetcher"
	"github.com/HulkLiu/WtTools/internal/model"
)

type Java666SearchEngine struct {
	E engine.Concurrent
}
type Java666 struct {
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

func (s *SearchManage) Search(url string) error {
	log.Printf("SearchEngine: Searching for %s", url)
	s.engine.Run(engine.Request{
		Url:        url,
		ParserFunc: s.parser.ClassPage,
	})
	return nil
}
func (j *Java666SearchEngine) Search(url string) error {

	data := Java666{
		ClassListRe: regexp.MustCompile(`<a target="_blank" href="(https:\/\/www.666php.com\/\d{1,}.html)" title=".*?" rel="bookmark">(.*?)<\/a>`),
		ClassPageRe: regexp.MustCompile(`(https:\/\/www.666php.com)`),
		UrlMap:      make(map[string]string),
		TitleRE:     regexp.MustCompile(`<h2>(.*)<\/h2>`),
		IdRe:        regexp.MustCompile(`com\/(.*?).html`),
		LimitItem:   5,
	}
	j.E.Run(engine.Request{
		Url:        url,
		ParserFunc: data.ClassPage,
	})
	return nil
}

func NewJava666() Parser {
	return &Java666{
		ClassListRe: regexp.MustCompile(`<a target="_blank" href="(https:\/\/www.666php.com\/\d{1,}.html)" title=".*?" rel="bookmark">(.*?)<\/a>`),
		ClassPageRe: regexp.MustCompile(`(https:\/\/www.666php.com)`),
		UrlMap:      make(map[string]string),
		TitleRE:     regexp.MustCompile(`<h2>(.*)<\/h2>`),
		IdRe:        regexp.MustCompile(`com\/(.*?).html`),
		LimitItem:   5,
	}
}

func (a *Java666) ClassPage(c []byte) engine.ParseResult {
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

func (a *Java666) ClassList(c []byte) engine.ParseResult {
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

func (a *Java666) ParseProfile(c []byte, url string, title string) engine.ParseResult {
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

func (a *Java666) extractString(c []byte, re *regexp.Regexp) string {
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
func (a *Java666) isUrlExists(url string) bool {
	if _, ok := a.UrlMap[url]; ok {
		return true
	}
	a.UrlMap[url] = url
	return false
}
