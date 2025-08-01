package service

import (
	"context"
	"time"

	v1 "todo-service/api/helloworld/v1"
	"todo-service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TodoService is a todo service.
type TodoService struct {
	v1.UnimplementedTodoServiceServer

	uc *biz.TodoUsecase
}

// NewTodoService new a todo service.
func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{uc: uc}
}

// CreateTodo implements helloworld.TodoServiceServer.
func (s *TodoService) CreateTodo(ctx context.Context, in *v1.CreateTodoRequest) (*v1.Todo, error) {
	todo, err := s.uc.CreateTodo(ctx, &biz.Todo{
		Title:       in.Title,
		Description: in.Description,
		Priority:    in.Priority,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Todo{
		Id:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Priority:    todo.Priority,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetTodo implements helloworld.TodoServiceServer.
func (s *TodoService) GetTodo(ctx context.Context, in *v1.GetTodoRequest) (*v1.Todo, error) {
	todo, err := s.uc.GetTodo(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.Todo{
		Id:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Priority:    todo.Priority,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateTodo implements helloworld.TodoServiceServer.
func (s *TodoService) UpdateTodo(ctx context.Context, in *v1.UpdateTodoRequest) (*v1.Todo, error) {
	todo, err := s.uc.UpdateTodo(ctx, &biz.Todo{
		ID:          in.Id,
		Title:       in.Title,
		Description: in.Description,
		Priority:    in.Priority,
		Status:      in.Status,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Todo{
		Id:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Priority:    todo.Priority,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteTodo implements helloworld.TodoServiceServer.
func (s *TodoService) DeleteTodo(ctx context.Context, in *v1.DeleteTodoRequest) (*emptypb.Empty, error) {
	err := s.uc.DeleteTodo(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ListTodos implements helloworld.TodoServiceServer.
func (s *TodoService) ListTodos(ctx context.Context, in *v1.ListTodosRequest) (*v1.ListTodosResponse, error) {
	todos, total, err := s.uc.ListTodos(ctx, in.Page, in.PageSize, in.Priority, in.Status)
	if err != nil {
		return nil, err
	}

	var respTodos []*v1.Todo
	for _, todo := range todos {
		respTodos = append(respTodos, &v1.Todo{
			Id:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			Status:      todo.Status,
			CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &v1.ListTodosResponse{
		Todos:    respTodos,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}