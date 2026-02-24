package models

import "time"

type Priority string
type Status string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"

	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Todo struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Tags        string     `json:"tags"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated,omitempty"`
}

type TodoEsDocument struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Tags        string     `json:"tags"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated"`
}

type CreateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Tags        string     `json:"tags"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type SearchTodoRequest struct {
	Query     string   `form:"query"`
	Status    Status   `form:"status"`
	Priority  Priority `form:"priority"`
	Tags      string   `form:"tags"`
	Page      int      `form:"page,default=1"`
	PageSize  int      `form:"page_size,default=20"`
	SortBy    string   `form:"sort_by,default=created_at"`
	SortOrder string   `form:"sort_order,default=desc"`
}

type SearchResponse struct {
	Total    int64  `json:"total"`
	Page     int    `json:"page"`
	Pagesize int    `json:"page_size"`
	Todos    []Todo `json:"todos"`
}
