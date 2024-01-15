package parser

import (
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/HulkLiu/WtTools/internal/fetcher"
	"github.com/HulkLiu/WtTools/internal/model"
	"log"
	"regexp"
	"strings"
	"time"
)

var (
//classListRe = regexp.MustCompile(`<a target="_blank" href="(https:\/\/www.666php.com\/\d{1,}.html)" title=".*?" rel="bookmark">(.*?)<\/a>`)
//
//classPageRe = regexp.MustCompile(`(https:\/\/www.666php.com)`)
//UrlMap      = make(map[string]string)
//
//TitleRE = regexp.MustCompile(`<h2>(.*)<\/h2>`)
//idRe    = regexp.MustCompile(`com\/(.*?).html`)
)

type Java666 struct {
	ClassListRe *regexp.Regexp

	ClassPageRe *regexp.Regexp
	UrlMap      map[string]string

	TitleRE *regexp.Regexp
	IdRe    *regexp.Regexp
}

func NewJava666() {

}

func (a *Java666) ClassPage(c []byte) engine.ParseResult {
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

		//log.Printf("Find %v ClassPageB url :%v,title:%v", nu, url, title)

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
		if _, ok := a.UrlMap[url]; ok {
			continue
		}
		a.UrlMap[url] = url

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
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func (a *Java666) ParseProfile(c []byte, url string, title string) engine.ParseResult {
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

func extractString(c []byte, re *regexp.Regexp) string {
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
