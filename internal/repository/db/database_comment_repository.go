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

type DatabaseCommentRepository struct {
	storage *storage.Storage[entity.Comment]
}

func NewDatabaseCommentRepository(storage *storage.Storage[entity.Comment]) *DatabaseCommentRepository {
	return &DatabaseCommentRepository{storage: storage}
}

func (d *DatabaseCommentRepository) Create(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	comment := &entity.Comment{}

	queryBuilder := squirrel.Insert("comment").
		Columns("post_id", "user_id", "content", "timestamp").
		Values(input.PostID, input.UserID, input.Content, input.Timestamp)
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return comment, fmt.Errorf("error building query: %w", err)
	}

	*comment, err = d.storage.ExecuteInsert(ctx, query, args...)
	if err != nil {
		return comment, fmt.Errorf("error inserting comment: %w", err)
	}

	return comment, nil
}

func (d *DatabaseCommentRepository) CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	comment := &entity.Comment{}

	queryBuilder := squirrel.Insert("comment").
		Columns("parent_id", "post_id", "user_id", "content", "timestamp").
		Values(input.ParentID, input.PostID, input.UserID, input.Content, input.Timestamp)
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return comment, fmt.Errorf("error building query: %w", err)
	}

	*comment, err = d.storage.ExecuteInsert(ctx, query, args...)
	if err != nil {
		return comment, fmt.Errorf("error inserting comment: %w", err)
	}

	return comment, nil
}

func (d *DatabaseCommentRepository) GetById(ctx context.Context, ID int,
	pagination entity.Pagination) (*entity.Comment, error) {
	comment := &entity.Comment{}

	queryBuilder := squirrel.Select("id", "parent_id", "post_id", "user_id", "content", "timestamp").
		From("comment").
		Where(squirrel.Eq{"id": ID})
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return comment, fmt.Errorf("error building query: %w", err)
	}

	err = d.storage.QueryRowContext(ctx, query, args).
		Scan(&comment.ID, &comment.ParentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return comment, fmt.Errorf("comment not found: %w", err)
		}
		return comment, fmt.Errorf("error fetching comment by ID: %w", err)
	}

	replies, err := d.getReplies(ctx, ID, pagination)
	if err != nil {
		return comment, fmt.Errorf("error fetching replies: %w", err)
	}

	comment.Replies = replies

	return comment, nil
}

func (d *DatabaseCommentRepository) getReplies(ctx context.Context, ID int,
	pagination entity.Pagination) ([]*entity.Comment, error) {
	var comments []*entity.Comment

	queryBuilder := squirrel.Select("id", "parent_id", "post_id", "user_id", "content", "timestamp").
		From("comment").
		OrderBy("timestamp").
		Where(squirrel.Eq{"parent_id": ID})

	if pagination.Limit != nil {
		queryBuilder = queryBuilder.Limit(uint64(*pagination.Limit))
	}
	if pagination.Offset != nil {
		queryBuilder = queryBuilder.Offset(uint64(*pagination.Offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return comments, fmt.Errorf("error building query: %w", err)
	}

	rows, err := d.storage.QueryContext(ctx, query, args)
	if err != nil {
		return comments, fmt.Errorf("error fetching replies comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		comment := &entity.Comment{}
		if err := rows.Scan(&comment.ID, &comment.ParentID, &comment.PostID, &comment.UserID, &comment.Content,
			&comment.Timestamp); err != nil {
			return comments, fmt.Errorf("error scanning comment: %w", err)
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return comments, fmt.Errorf("error iterating posts: %w", err)
	}

	return comments, nil
}

func (d *DatabaseCommentRepository) GetAll(ctx context.Context, filter entity.CommentFilter,
	pagination entity.Pagination) ([]*entity.Comment, error) {
	var comments []*entity.Comment

	queryBuilder := squirrel.Select("id", "parent_id", "post_id", "user_id", "content", "timestamp").
		From("comment").
		OrderBy("timestamp")

	if filter.UserID != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": filter.UserID})
	}

	if pagination.Limit != nil {
		queryBuilder = queryBuilder.Limit(uint64(*pagination.Limit))
	}
	if pagination.Offset != nil {
		queryBuilder = queryBuilder.Offset(uint64(*pagination.Offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return comments, fmt.Errorf("error building SQL query: %w", err)
	}

	rows, err := d.storage.QueryContext(ctx, query, args)
	if err != nil {
		return comments, fmt.Errorf("error fetching posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		comment := &entity.Comment{}
		if err := rows.Scan(&comment.ID, &comment.ParentID, &comment.PostID, &comment.UserID,
			&comment.Content, &comment.Timestamp); err != nil {
			return comments, fmt.Errorf("error scanning comment: %w", err)
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return comments, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

func (d *DatabaseCommentRepository) Delete(ctx context.Context, ID int) error {
	queryBuilder := squirrel.Delete("comment").
		Where(squirrel.Eq{"post_id": ID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query: %w", err)
	}

	err = d.storage.ExecuteDelete(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}
