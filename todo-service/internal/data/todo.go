package data

import (
	"context"
	"database/sql"
	"time"

	"todo-service/api/helloworld/v1"
	"todo-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type todoRepo struct {
	data *Data
	log  *log.Helper
}

// NewTodoRepo .
func NewTodoRepo(data *Data, logger log.Logger) biz.TodoRepo {
	return &todoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *todoRepo) initTable(ctx context.Context) error {
	_, err := r.data.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS todos (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			priority INT NOT NULL DEFAULT 2,
			status INT NOT NULL DEFAULT 1,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

func (r *todoRepo) Save(ctx context.Context, todo *biz.Todo) (*biz.Todo, error) {
	// 确保表存在
	if err := r.initTable(ctx); err != nil {
		r.log.WithContext(ctx).Errorf("Failed to init table: %v", err)
		return nil, err
	}

	result, err := r.data.db.ExecContext(ctx, `
		INSERT INTO todos (title, description, priority, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, todo.Title, todo.Description, todo.Priority, todo.Status, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save todo: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.log.WithContext(ctx).Errorf("Failed to get last insert id: %v", err)
		return nil, err
	}
	todo.ID = id

	return todo, nil
}

func (r *todoRepo) Update(ctx context.Context, todo *biz.Todo) (*biz.Todo, error) {
	_, err := r.data.db.ExecContext(ctx, `
		UPDATE todos
		SET title = ?, description = ?, priority = ?, status = ?, updated_at = ?
		WHERE id = ?
	`, todo.Title, todo.Description, todo.Priority, todo.Status, todo.UpdatedAt, todo.ID)
	if err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update todo: %v", err)
		return nil, err
	}

	return todo, nil
}

func (r *todoRepo) FindByID(ctx context.Context, id int64) (*biz.Todo, error) {
	var todo biz.Todo
	var createdAt, updatedAt time.Time

	err := r.data.db.QueryRowContext(ctx, `
		SELECT id, title, description, priority, status, created_at, updated_at
		FROM todos
		WHERE id = ?
	`, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Priority, &todo.Status,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.WithContext(ctx).Infof("Todo not found: %d", id)
			return nil, err
		}
		r.log.WithContext(ctx).Errorf("Failed to find todo: %v", err)
		return nil, err
	}

	todo.CreatedAt = createdAt
	todo.UpdatedAt = updatedAt

	return &todo, nil
}

func (r *todoRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.data.db.ExecContext(ctx, `
		DELETE FROM todos
		WHERE id = ?
	`, id)
	if err != nil {
		r.log.WithContext(ctx).Errorf("Failed to delete todo: %v", err)
		return err
	}

	return nil
}

func (r *todoRepo) List(ctx context.Context, page, pageSize int32, priority v1.Priority, status v1.Status) ([]*biz.Todo, int32, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := `
		SELECT id, title, description, priority, status, created_at, updated_at
		FROM todos
		WHERE 1=1
	`
	countQuery := `
		SELECT COUNT(*)
		FROM todos
		WHERE 1=1
	`

	params := []interface{}{}

	// 添加优先级过滤
	if priority != v1.Priority_PRIORITY_UNSPECIFIED {
		query += " AND priority = ?"
		countQuery += " AND priority = ?"
		params = append(params, priority)
	}

	// 添加状态过滤
	if status != v1.Status_STATUS_UNSPECIFIED {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		params = append(params, status)
	}

	// 添加分页
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	params = append(params, pageSize, offset)

	// 执行查询
	rows, err := r.data.db.QueryContext(ctx, query, params...)
	if err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list todos: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	// 处理结果
	todos := []*biz.Todo{}
	for rows.Next() {
		var todo biz.Todo
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Priority, &todo.Status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			r.log.WithContext(ctx).Errorf("Failed to scan todo: %v", err)
			return nil, 0, err
		}

		todo.CreatedAt = createdAt
		todo.UpdatedAt = updatedAt
		todos = append(todos, &todo)
	}

	// 获取总记录数
	var total int32
	countParams := params[:len(params)-2] // 移除分页参数
	err = r.data.db.QueryRowContext(ctx, countQuery, countParams...).Scan(&total)
	if err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count todos: %v", err)
		return nil, 0, err
	}

	return todos, total, nil
}