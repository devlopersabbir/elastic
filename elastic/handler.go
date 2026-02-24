package elastic

import (
	"bytes"
	"context"
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

func DeleteTodo(id uint) {
	if EsClient == nil {
		return
	}
	res, err := EsClient.Delete("todos", strconv.Itoa(int(id)))
	if err != nil {
		log.Println("delete request failed:", err)
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Println("delete error:", res.String())
	}
}

func Search(ctx context.Context, req *models.SearchTodoRequest) ([]uint, int64, error) {
	if EsClient == nil {
		return nil, 0, nil
	}

	mustArr := []interface{}{}

	if req.Query != "" {
		mustArr = append(mustArr, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"title^3", "description"},
			},
		})
	}
	if req.Status != "" {
		mustArr = append(mustArr, map[string]interface{}{
			"term": map[string]interface{}{
				"status": req.Status,
			},
		})
	}
	if req.Priority != "" {
		mustArr = append(mustArr, map[string]interface{}{
			"term": map[string]interface{}{
				"priority": req.Priority,
			},
		})
	}
	if req.Tags != "" {
		mustArr = append(mustArr, map[string]interface{}{
			"match": map[string]interface{}{
				"tags": req.Tags,
			},
		})
	}

	var query interface{}
	if len(mustArr) == 0 {
		query = map[string]interface{}{"match_all": map[string]interface{}{}}
	} else {
		query = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustArr,
			},
		}
	}

	queryMap := map[string]interface{}{
		"query": query,
		"from":  (req.Page - 1) * req.PageSize,
		"size":  req.PageSize,
	}

	queryBytes, _ := json.Marshal(queryMap)
	res, err := EsClient.Search(
		EsClient.Search.WithContext(ctx),
		EsClient.Search.WithIndex("todos"),
		EsClient.Search.WithBody(bytes.NewReader(queryBytes)),
		EsClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println("search error:", res.String())
		return nil, 0, nil
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}

	hitsMap, ok := r["hits"].(map[string]interface{})
	if !ok {
		return nil, 0, nil
	}

	hits, ok := hitsMap["hits"].([]interface{})
	if !ok {
		return nil, 0, nil
	}

	var total int64
	if totalMap, ok := hitsMap["total"].(map[string]interface{}); ok {
		if val, ok := totalMap["value"].(float64); ok {
			total = int64(val)
		}
	}

	var ids []uint
	for _, hit := range hits {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			if source, ok := hitMap["_source"].(map[string]interface{}); ok {
				if idFloat, ok := source["id"].(float64); ok {
					ids = append(ids, uint(idFloat))
				}
			}
		}
	}

	return ids, total, nil
}
