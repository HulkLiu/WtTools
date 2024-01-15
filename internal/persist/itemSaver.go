package persist

import (
	"context"

	"fmt"

	"github.com/HulkLiu/WtTools/internal/engine"
	"github.com/olivere/elastic/v7"
)

// http://localhost:9200/php6661/_search
var err error

func ItemSaver(index string, client *elastic.Client) (chan engine.Item, error) {
	//client, err := elastic.NewClient(elastic.SetSniff(false))
	//if err != nil {
	//	return nil, err
	//}

	out := make(chan engine.Item)
	// count := 0
	go func() {
		for {
			item := <-out
			//fmt.Printf("%v -> got data %v\n", count, item)
			// if config.IsSAVE {
			// 	if err = save(item, client, index); err != nil {
			// 		log.Printf("save false data %+v\n", item.Url)
			// 	} else {
			// 		count++
			// 	}
			// }
			fmt.Printf("Has successfully obtained data for %v items\n", item.Url)

		}
	}()

	return out, nil
}
func save(item engine.Item, client *elastic.Client, index string) (err error) {

	indexServer := client.Index().
		Index(index).
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
