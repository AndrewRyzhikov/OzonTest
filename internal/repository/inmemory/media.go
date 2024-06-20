package inmemory

import (
	"context"
	"sync"

	"OzonTest/internal/entity"
)

type posts map[int]*post
type Media struct {
	posts        posts
	keyGenerator int

	sync.RWMutex
}

func NewMedia() *Media {
	return &Media{posts: make(posts)}
}

func (media *Media) CreatePost(ctx context.Context, input entity.Post) (*entity.Post, error) {
	media.Lock()
	defer media.Unlock()

	p, ok := media.posts[input.PostID]
	if ok {
		return &entity.Post{}, entity.ErrAlreadyExists
	}

	p.insert(input)

	return toEntityPost(p), nil
}

func (media *Media) GetByIdPost(ctx context.Context, ID int) (*entity.Post, error) {
	media.RLock()
	defer media.RUnlock()

	p, ok := media.posts[ID]
	if !ok {
		return &entity.Post{}, entity.ErrNotFound
	}

	return toEntityPost(p), nil
}

func (media *Media) GetAllPosts(ctx context.Context, filter entity.PostFilter, pagination entity.Pagination) ([]*entity.Post, error) {
	media.RLock()
	defer media.RUnlock()
	posts := make([]*entity.Post, 0, len(media.posts))

	for _, p := range media.posts {
		posts = append(posts, toEntityPost(p))
	}
	return posts, nil
}

func (media *Media) DeletePost(ctx context.Context, ID int) error {
	//TODO implement me
	panic("implement me")
}

func (media *Media) SwitchCommentsState(ctx context.Context, ID int, state bool) (*entity.Post, error) {
	//TODO implement me
	panic("implement me")
}

func NewCommentInMemory() *Media {
	return &Media{posts: make(posts), keyGenerator: 1}
}

func (media *Media) CreateComment(_ context.Context, input entity.Comment) (*entity.Comment, error) {
	media.Lock()
	defer media.Unlock()

	input.ID = media.keyGenerator
	//media.comments[media.keyGenerator] = &input

	media.keyGenerator++

	return &input, nil
}

func (media *Media) CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	media.Lock()

	//if _, ok := media.comments[input.ParentID]; !ok {
	//	return &entity.Comment{}, errors.New("parent ID not found")
	//}

	media.Unlock()

	return media.CreateComment(ctx, input)
}

func (media *Media) GetCommentById(_ context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (media *Media) GetAllComments(_ context.Context, filter entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (media *Media) DeleteComment(ctx context.Context, ID int) error {
	//TODO implement me
	panic("implement me")
}
