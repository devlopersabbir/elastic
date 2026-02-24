package elastic

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"

	"github.com/devloeprsabbir/go-elasticsearch/models"
)

func IndexTodo(t *models.Todo) {
	jsonBody, err := json.Marshal(t)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}

	res, err := EsClient.Index(
		"todos",
		bytes.NewReader(jsonBody),
		EsClient.Index.WithDocumentID(strconv.Itoa(int(t.ID))), // VERY IMPORTANT
		EsClient.Index.WithRefresh("true"),                     // force refresh so you can see immediately
	)
	if err != nil {
		log.Println("index request failed:", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println("indexing error:", res.String())
		return
	}

	log.Println("indexed todo successfully:", t.ID)
}
