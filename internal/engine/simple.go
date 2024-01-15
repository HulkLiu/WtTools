package engine

import (
	"fmt"
)

type SimpleEngine struct {
}

func (e SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		res, err := Worker(r)
		if err != nil {
			fmt.Printf("fetcher false Url:%v", r.Url)

			continue
		}
		requests = append(requests, res.Requests...)

		//for k, item := range res.Items {
		//	fmt.Printf("%v -> url : %v ,title -> %s\n", k+1, r.Url, item)
		//}
		//break
	}

}
