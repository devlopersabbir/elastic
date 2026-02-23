package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *TodoRepository) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.WithContext(ctx).First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) GetByIDs(ctx context.Context, ids []uint) ([]models.Todo, error) {
	if len(ids) == 0 {
		return []models.Todo{}, nil
	}
	var todos []models.Todo
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&todos).Error
	if err != nil {
		return nil, err
	}

	// Preserve the order from Elasticsearch results
	todoMap := make(map[uint]models.Todo, len(todos))
	for _, t := range todos {
		todoMap[t.ID] = t
	}
	ordered := make([]models.Todo, 0, len(ids))
	for _, id := range ids {
		if t, ok := todoMap[id]; ok {
			ordered = append(ordered, t)
		}
	}
	return ordered, nil
}

func (r *TodoRepository) Update(ctx context.Context, id uint, req *models.UpdateTodoRequest) (*models.Todo, error) {
	var todo models.Todo
	if err := r.db.WithContext(ctx).First(&todo, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Priority != "" {
		updates["priority"] = req.Priority
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.DueDate != nil {
		updates["due_date"] = req.DueDate
	}

	if err := r.db.WithContext(ctx).Model(&todo).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Todo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}
	return nil
}

func (r *TodoRepository) GetAll(ctx context.Context, page, pageSize int) ([]models.Todo, int64, error) {
	var todos []models.Todo
	var total int64

	r.db.WithContext(ctx).Model(&models.Todo{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&todos).Error

	return todos, total, err
}

// BulkInsert for seeding/testing with millions of todos
func (r *TodoRepository) BulkInsert(ctx context.Context, todos []models.Todo) error {
	return r.db.WithContext(ctx).CreateInBatches(todos, 1000).Error
}