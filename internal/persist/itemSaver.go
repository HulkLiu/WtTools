package persist

import (
	"context"

	"fmt"
	"github.com/HulkLiu/WtTools/internal/config"
	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/olivere/elastic/v7"
)

//http://localhost:9200/php6661/_search

func ItemSaver(index string) (chan engine.Item, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	count := 0
	go func() {
		for {
			item := <-out
			//fmt.Printf("%v -> got data %v\n", count, item)
			if config.IsSAVE {
				if err = save(item, client, index); err != nil {
					fmt.Printf("save false data %+v\n", item)
				} else {
					count++
				}
			}
			fmt.Printf("Has successfully obtained data for %v items\n", count)

		}
	}()

	return out, nil
}
func save(item engine.Item, client *elastic.Client, index string) (err error) {

	indexServer := client.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexServer.Id(item.Id)
	}
	_, err = indexServer.
		Do(context.Background())

	if err != nil {
		return err

	}
	return nil

}
