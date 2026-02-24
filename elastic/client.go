package elastic

import (
	"log"

	config "github.com/devloeprsabbir/go-elasticsearch/config"
	"github.com/elastic/go-elasticsearch/v9"
)

var EsClient *elasticsearch.Client

func ElasticClient(cfg *config.Config) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchUrl},
	})
	if err != nil {
		log.Fatalf("failed to connect to elasticsearch: %v", err)
	}

	// create index
	CreateTodoIndex(es)
	EsClient = es
}

func CreateTodoIndex(es *elasticsearch.Client) {
	res, err := es.Indices.Create("todos")
	if err != nil {
		log.Fatalf("failed creating index: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println("index may already exist:", res.String())
	}
}
