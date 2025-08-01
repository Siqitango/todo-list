package biz

import (
	"context"
	"time"

	v1 "todo-service/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrTodoNotFound is todo not found.
	ErrTodoNotFound = errors.NotFound(v1.ErrorReason_TODO_NOT_FOUND.String(), "todo not found")
)

// Todo is a Todo model.
type Todo struct {
	ID          int64
	Title       string
	Description string
	Priority    v1.Priority
	Status      v1.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TodoRepo is a Todo repo.
type TodoRepo interface {
	Save(context.Context, *Todo) (*Todo, error)
	Update(context.Context, *Todo) (*Todo, error)
	FindByID(context.Context, int64) (*Todo, error)
	Delete(context.Context, int64) error
	List(context.Context, int32, int32, v1.Priority, v1.Status) ([]*Todo, int32, error)
}

// TodoUsecase is a Todo usecase.
type TodoUsecase struct {
	repo TodoRepo
	log  *log.Helper
}

// NewTodoUsecase new a Todo usecase.
func NewTodoUsecase(repo TodoRepo, logger log.Logger) *TodoUsecase {
	return &TodoUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateTodo creates a Todo, and returns the new Todo.
func (uc *TodoUsecase) CreateTodo(ctx context.Context, todo *Todo) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("CreateTodo: %v", todo.Title)
	todo.Status = v1.Status_STATUS_PENDING
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	return uc.repo.Save(ctx, todo)
}

// GetTodo gets a Todo by id.
func (uc *TodoUsecase) GetTodo(ctx context.Context, id int64) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("GetTodo: %d", id)
	todo, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTodoNotFound
	}
	return todo, nil
}

// UpdateTodo updates a Todo.
func (uc *TodoUsecase) UpdateTodo(ctx context.Context, todo *Todo) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("UpdateTodo: %d", todo.ID)
	existing, err := uc.repo.FindByID(ctx, todo.ID)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	if todo.Title != "" {
		existing.Title = todo.Title
	}
	if todo.Description != "" {
		existing.Description = todo.Description
	}
	if todo.Priority != v1.Priority_PRIORITY_UNSPECIFIED {
		existing.Priority = todo.Priority
	}
	if todo.Status != v1.Status_STATUS_UNSPECIFIED {
		existing.Status = todo.Status
	}
	existing.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, existing)
}

// DeleteTodo deletes a Todo.
func (uc *TodoUsecase) DeleteTodo(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("DeleteTodo: %d", id)
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return ErrTodoNotFound
	}
	return uc.repo.Delete(ctx, id)
}

// ListTodos lists Todos.
func (uc *TodoUsecase) ListTodos(ctx context.Context, page, pageSize int32, priority v1.Priority, status v1.Status) ([]*Todo, int32, error) {
	uc.log.WithContext(ctx).Infof("ListTodos: page=%d, pageSize=%d, priority=%v, status=%v", page, pageSize, priority, status)
	return uc.repo.List(ctx, page, pageSize, priority, status)
}