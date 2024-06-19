package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"OzonTest/internal/entity"
	"OzonTest/internal/repository/db/storage"
)

type DatabasePostRepository struct {
	storage *storage.Storage[entity.Post]
}

func NewDatabasePostRepository(storage *storage.Storage[entity.Post]) *DatabasePostRepository {
	return &DatabasePostRepository{storage: storage}
}

func (d *DatabasePostRepository) Create(ctx context.Context, input entity.Post) (*entity.Post, error) {
	post := &entity.Post{}

	queryBuilder := squirrel.Insert("post").
		Columns("user_id", "content", "is_open", "timestamp").
		Values(input.UserID, input.Content, input.IsOpen, input.Timestamp)
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return post, fmt.Errorf("error building query: %w", err)
	}

	*post, err = d.storage.ExecuteInsert(ctx, query, args...)
	if err != nil {
		return post, fmt.Errorf("error inserting post: %w", err)
	}

	return post, nil
}

func (d *DatabasePostRepository) GetById(ctx context.Context, ID int) (*entity.Post, error) {
	post := &entity.Post{}

	queryBuilder := squirrel.Select("id", "user_id", "content", "is_open", "timestamp").
		From("post").
		Where(squirrel.Eq{"id": ID})
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return post, fmt.Errorf("error building query: %w", err)
	}

	err = d.storage.QueryRowContext(ctx, query, args).Scan(&post.PostID, &post.UserID, &post.Content,
		&post.IsOpen, &post.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return post, fmt.Errorf("post not found: %w", err)
		}
		return post, fmt.Errorf("error fetching post by ID: %w", err)
	}
	return post, nil
}

func (d *DatabasePostRepository) GetAll(ctx context.Context, filter entity.PostFilter,
	pagination entity.Pagination) ([]*entity.Post, error) {
	var posts []*entity.Post

	queryBuilder := squirrel.Select("id", "user_id", "content", "is_open", "timestamp").
		OrderBy("timestamp").
		From("post")

	if filter.UserID != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": filter.UserID})
	}

	queryBuilder = queryBuilder.Limit(uint64(*pagination.Limit)).Offset(uint64(*pagination.Offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return posts, fmt.Errorf("error building SQL query: %w", err)
	}

	rows, err := d.storage.QueryContext(ctx, query, args)
	if err != nil {
		return posts, fmt.Errorf("error fetching posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := &entity.Post{}
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.IsOpen, &post.Timestamp); err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

func (d *DatabasePostRepository) Delete(ctx context.Context, ID int) error {
	queryBuilder := squirrel.Delete("post").
		Where(squirrel.Eq{"id": ID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query: %w", err)
	}

	err = d.storage.ExecuteDelete(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (d *DatabasePostRepository) SwitchCommentsState(ctx context.Context, ID int, flag bool) (*entity.Post, error) {
	post := &entity.Post{}

	queryBuilder := squirrel.Update("post").
		Set("is_open", flag).
		Where(squirrel.Eq{"id": ID}).
		Suffix("RETURNING *")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return post, fmt.Errorf("failed to build query: %w", err)
	}

	*post, err = d.storage.ExecuteUpdate(ctx, query, args...)
	if err != nil {
		return post, err
	}

	return post, nil
}
