package elastic

import (
	"log"

	config "github.com/devloeprsabbir/go-elasticsearch/config"
	"github.com/elastic/go-elasticsearch/v9"
)

var EsClient *elasticsearch.Client

func ElasticClient (cfg *config.Config) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchUrl},
	})
	if err != nil {
		log.Fatalf("failed to connect to elasticsearch: %v", err)
	}

	EsClient = es 
}