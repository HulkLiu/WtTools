package engine

import (
	"github.com/HulkLiu/WtTools/internal/fetcher"
	"log"
)

func Worker(r Request) (ParseResult, error) {
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error "+"fetching url %s: %v", r.Url, err)

		return ParseResult{}, err
	}
	return r.ParserFunc(body), nil
}
