package handler

import (
	"net/http"
	"strconv"

	"github.com/devloeprsabbir/go-elasticsearch/elastic"
	"github.com/devloeprsabbir/go-elasticsearch/models"
	"github.com/devloeprsabbir/go-elasticsearch/repository"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	repo *repository.TodoRepository
	es   *elasticsearch.Client
}

func NewTodoHandler(repo *repository.TodoRepository, es *elasticsearch.Client) *TodoHandler {
	return &TodoHandler{
		repo: repo,
		es:   es,
	}
}

// CreateTodo godoc
// POST /todos
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Priority == "" {
		req.Priority = models.PriorityMedium
	}

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Tags:        req.Tags,
		DueDate:     req.DueDate,
		Status:      models.StatusPending,
	}

	if err := h.repo.Create(c.Request.Context(), todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create todo"})
		return
	}

	// Index in Elasticsearch asynchronously
	go elastic.IndexTodo(todo)

	c.JSON(http.StatusCreated, todo)
}

// GetTodo godoc
// GET /todos/:id
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// // UpdateTodo godoc
// // PUT /todos/:id
// func (h *TodoHandler) UpdateTodo(c *gin.Context) {
// 	id, err := parseUintParam(c, "id")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
// 		return
// 	}

// 	var req models.UpdateTodoRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	todo, err := h.repo.Update(c.Request.Context(), id, &req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo"})
// 		return
// 	}

// 	// Re-index in Elasticsearch
// 	go func() {
// 		if err := h.es.IndexTodo(context.Background(), todo); err != nil {
// 			_ = err
// 		}
// 	}()

// 	c.JSON(http.StatusOK, todo)
// }

// // DeleteTodo godoc
// // DELETE /todos/:id
// func (h *TodoHandler) DeleteTodo(c *gin.Context) {
// 	id, err := parseUintParam(c, "id")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
// 		return
// 	}

// 	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
// 		return
// 	}

// 	// Remove from Elasticsearch
// 	go func() {
// 		if err := h.es.DeleteTodo(context.Background(), id); err != nil {
// 			_ = err
// 		}
// 	}()

// 	c.JSON(http.StatusOK, gin.H{"message": "todo deleted successfully"})
// }

// // ListTodos godoc
// // GET /todos - simple paginated list without search
// func (h *TodoHandler) ListTodos(c *gin.Context) {
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 || pageSize > 100 {
// 		pageSize = 20
// 	}

// 	todos, total, err := h.repo.GetAll(c.Request.Context(), page, pageSize)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list todos"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.SearchResponse{
// 		Total:    total,
// 		Page:     page,
// 		PageSize: pageSize,
// 		Todos:    todos,
// 	})
// }

// // SearchTodos godoc
// // GET /todos/search - full-text search via Elasticsearch
// func (h *TodoHandler) SearchTodos(c *gin.Context) {
// 	var req models.SearchRequest
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if req.Page < 1 {
// 		req.Page = 1
// 	}
// 	if req.PageSize < 1 || req.PageSize > 100 {
// 		req.PageSize = 20
// 	}

// 	// Search in Elasticsearch to get IDs
// 	ids, total, err := h.es.Search(c.Request.Context(), &req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
// 		return
// 	}

// 	// Fetch full documents from PostgreSQL in the same order
// 	todos, err := h.repo.GetByIDs(c.Request.Context(), ids)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch results"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.SearchResponse{
// 		Total:    total,
// 		Page:     req.Page,
// 		PageSize: req.PageSize,
// 		Todos:    todos,
// 	})
// }

// // SeedTodos godoc
// // POST /todos/seed - generate N sample todos for testing at scale
// func (h *TodoHandler) SeedTodos(c *gin.Context) {
// 	var body struct {
// 		Count int `json:"count" binding:"required,min=1,max=100000"`
// 	}
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	go func() {
// 		generateAndIndexTodos(h, body.Count)
// 	}()

// 	c.JSON(http.StatusAccepted, gin.H{
// 		"message": "Seeding started in background",
// 		"count":   body.Count,
// 	})
// }

// func generateAndIndexTodos(h *TodoHandler, count int) {
// 	titles := []string{
// 		"Buy groceries", "Write report", "Call dentist", "Fix bug in auth module",
// 		"Review pull request", "Update documentation", "Deploy to production",
// 		"Send invoice to client", "Prepare meeting agenda", "Refactor database queries",
// 		"Write unit tests", "Schedule team meeting", "Update dependencies",
// 		"Analyze performance metrics", "Design new feature", "Onboard new employee",
// 	}
// 	priorities := []models.Priority{models.PriorityLow, models.PriorityMedium, models.PriorityHigh}
// 	statuses := []models.Status{models.StatusPending, models.StatusInProgress, models.StatusDone}
// 	tags := []string{"work", "personal", "urgent", "backend", "frontend", "devops", "meeting", "finance"}

// 	batchSize := 1000
// 	for i := 0; i < count; i += batchSize {
// 		end := i + batchSize
// 		if end > count {
// 			end = count
// 		}

// 		todos := make([]models.Todo, 0, end-i)
// 		for j := i; j < end; j++ {
// 			todos = append(todos, models.Todo{
// 				Title:       titles[j%len(titles)] + " #" + strconv.Itoa(j+1),
// 				Description: "Auto-generated todo number " + strconv.Itoa(j+1),
// 				Status:      statuses[j%len(statuses)],
// 				Priority:    priorities[j%len(priorities)],
// 				Tags:        tags[j%len(tags)] + " " + tags[(j+1)%len(tags)],
// 			})
// 		}

// 		if err := h.repo.BulkInsert(context.Background(), todos); err != nil {
// 			return
// 		}

// 		// Re-fetch to get IDs assigned by DB
// 		savedTodos, _, err := h.repo.GetAll(context.Background(), (i/batchSize)+1, batchSize)
// 		if err != nil {
// 			return
// 		}

// 		if err := h.es.BulkIndex(context.Background(), savedTodos); err != nil {
// 			return
// 		}
// 	}
// }

func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
