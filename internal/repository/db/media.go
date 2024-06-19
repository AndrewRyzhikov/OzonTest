package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"OzonTest/internal/entity"
)

type DatabaseMedia struct {
	commentStorage *Storage[entity.Comment]
	postStorage    *Storage[entity.Post]
}

func NewDatabaseMedia(commentStorage *Storage[entity.Comment], postStorage *Storage[entity.Post]) *DatabaseMedia {
	return &DatabaseMedia{commentStorage: commentStorage, postStorage: postStorage}
}

func (d *DatabaseMedia) CreateComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	comment := &entity.Comment{}

	queryBuilder := squirrel.Insert("comment").
		Columns("post_id", "user_id", "content", "timestamp").
		Values(input.PostID, input.UserID, input.Content, input.Timestamp)
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return comment, fmt.Errorf("error building query: %w", err)
	}

	*comment, err = d.commentStorage.Insert(ctx, query, args...)
	if err != nil {
		return comment, fmt.Errorf("error inserting comment: %w", err)
	}

	return comment, nil
}

func (d *DatabaseMedia) CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	comment := &entity.Comment{}

	queryBuilder := squirrel.Insert("comment").
		Columns("parent_id", "post_id", "user_id", "content", "timestamp").
		Values(input.ParentID, input.PostID, input.UserID, input.Content, input.Timestamp)
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return comment, fmt.Errorf("error building query: %w", err)
	}

	*comment, err = d.commentStorage.Insert(ctx, query, args...)
	if err != nil {
		return comment, fmt.Errorf("error inserting comment: %w", err)
	}

	return comment, nil
}

func (d *DatabaseMedia) GetCommentById(ctx context.Context, ID int,
	pagination entity.Pagination) (*entity.Comment, error) {
	comment := &entity.Comment{}

	queryBuilder := squirrel.Select("id", "parent_id", "post_id", "user_id", "content", "timestamp").
		From("comment").
		Where(squirrel.Eq{"id": ID})
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return comment, fmt.Errorf("error building query: %w", err)
	}

	err = d.commentStorage.QueryRowContext(ctx, query, args).
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

func (d *DatabaseMedia) getReplies(ctx context.Context, ID int,
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

	rows, err := d.commentStorage.QueryContext(ctx, query, args)
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

func (d *DatabaseMedia) GetAllComments(ctx context.Context, filter entity.CommentFilter,
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

	rows, err := d.commentStorage.QueryContext(ctx, query, args)
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

func (d *DatabaseMedia) DeleteComment(ctx context.Context, ID int) error {
	queryBuilder := squirrel.Delete("comment").
		Where(squirrel.Eq{"post_id": ID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query: %w", err)
	}

	err = d.commentStorage.Delete(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (d *DatabaseMedia) CreatePost(ctx context.Context, input entity.Post) (*entity.Post, error) {
	post := &entity.Post{}

	queryBuilder := squirrel.Insert("post").
		Columns("user_id", "content", "is_open", "timestamp").
		Values(input.UserID, input.Content, input.IsOpen, input.Timestamp)
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return post, fmt.Errorf("error building query: %w", err)
	}

	*post, err = d.postStorage.Insert(ctx, query, args...)
	if err != nil {
		return post, fmt.Errorf("error inserting post: %w", err)
	}

	return post, nil
}

func (d *DatabaseMedia) GetByIdPost(ctx context.Context, ID int) (*entity.Post, error) {
	post := &entity.Post{}

	queryBuilder := squirrel.Select("id", "user_id", "content", "is_open", "timestamp").
		From("post").
		Where(squirrel.Eq{"id": ID})
	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return post, fmt.Errorf("error building query: %w", err)
	}

	err = d.postStorage.QueryRowContext(ctx, query, args).Scan(&post.PostID, &post.UserID, &post.Content,
		&post.IsOpen, &post.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return post, fmt.Errorf("post not found: %w", err)
		}
		return post, fmt.Errorf("error fetching post by ID: %w", err)
	}
	return post, nil
}

func (d *DatabaseMedia) GetAllPosts(ctx context.Context, filter entity.PostFilter,
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

	rows, err := d.postStorage.QueryContext(ctx, query, args)
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

func (d *DatabaseMedia) DeletePost(ctx context.Context, ID int) error {
	queryBuilder := squirrel.Delete("post").
		Where(squirrel.Eq{"id": ID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query: %w", err)
	}

	err = d.postStorage.Delete(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (d *DatabaseMedia) SwitchCommentsState(ctx context.Context, ID int, flag bool) (*entity.Post, error) {
	post := &entity.Post{}

	queryBuilder := squirrel.Update("post").
		Set("is_open", flag).
		Where(squirrel.Eq{"id": ID}).
		Suffix("RETURNING *")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return post, fmt.Errorf("failed to build query: %w", err)
	}

	*post, err = d.postStorage.ExecuteUpdate(ctx, query, args...)
	if err != nil {
		return post, err
	}

	return post, nil
}
