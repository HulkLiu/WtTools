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

type Xue600SearchEngine struct {
	E engine.Concurrent
}

type Xue600 struct {
	ClassListRe     *regexp.Regexp
	ClassPageRe     *regexp.Regexp
	UrlMap          map[string]string
	TitleRE         *regexp.Regexp
	IdRe            *regexp.Regexp
	LanguageTypeSel string
	CatalogSel      string
	LastTimeSel     string
	LimitItem       int
}

func (x *Xue600SearchEngine) Search(url string) error {
	//url := "https://www.600xue.com/?s=" + keyword
	log.Printf("Xue600SearchEngine: Searching for %s", url)
	// Xue600 specific parsing logic
	data := Xue600{
		ClassListRe:     regexp.MustCompile(`<a target="_blank" href="(.*?)" title=".*?" rel="bookmark">(.*?)<\/a>`),
		ClassPageRe:     regexp.MustCompile(`(https:\/\/www.600xue.com)`),
		UrlMap:          make(map[string]string),
		TitleRE:         regexp.MustCompile(`<h2>(.*)<\/h2>`),
		IdRe:            regexp.MustCompile(`com\/(.*?)\/`),
		LanguageTypeSel: "body > div > div.site-content > div > div.breadcrumbs > a:nth-child(3)",
		LimitItem:       5,
	}
	x.E.Run(engine.Request{
		Url:        url,
		ParserFunc: data.ClassPage,
	})
	return nil
}

func (a *Xue600) ClassPage(c []byte) engine.ParseResult {
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

		//fmt.Printf("ClassList => ClassPageA url %v\n", url)
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
		p.Title = a.extractString(c, a.TitleRE)
	}

	p.LanguageType = fetcher.GetTextByQuery(string(c), a.LanguageTypeSel)

	p.Catalog = fetcher.GetTextByQuery(string(c), "#post-"+a.extractString([]byte(url), regexp.MustCompile(`com\/(.*?)\/`))+" > div:nth-child(2) > div > div.entry-content.u-text-format.u-clearfix > p")

	p.LastTime = fetcher.GetTextByQuery(string(c), "#post-"+a.extractString([]byte(url), regexp.MustCompile(`com\/(.*?)\/`))+" > div:nth-child(1) > div > header > div > span.meta-date > a > time")

	p.CreateAt = time.Now().Format(config.DateFormat)
	p.Status = "待确认"

	//log.Printf("%+v", p)
	//utils.DownFile(fmt.Sprintf("%+v", p))

	res := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "Xue600",
				Id:      a.extractString([]byte(url), a.IdRe),
				Payload: p,
			},
		},
	}

	//fmt.Printf("%+v\n", res)
	return res
}
func (a *Xue600) extractString(c []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(c)

	if len(match) > 1 {
		return string(match[1])
	} else {
		return ""
	}
}
