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
	return &Media{posts: make(posts), keyGenerator: 1}
}

func (media *Media) CreatePost(_ context.Context, input entity.Post) (*entity.Post, error) {
	media.Lock()
	defer media.Unlock()

	_, ok := media.posts[input.ID]
	if ok {
		return &entity.Post{}, entity.ErrAlreadyExists("post", input.ID)
	}

	newPost := newPost(media.keyGenerator, input)
	media.posts[newPost.PostID] = newPost
	media.keyGenerator++

	return toEntityPost(newPost), nil
}

func (media *Media) GetByIdPost(_ context.Context, ID int) (*entity.Post, error) {
	media.RLock()
	defer media.RUnlock()

	p, ok := media.posts[ID]
	if !ok {
		return &entity.Post{}, entity.ErrNotFound("post", ID)
	}

	return toEntityPost(p), nil
}

func (media *Media) GetAllPosts(_ context.Context, filter *entity.PostFilter, pagination entity.Pagination) ([]*entity.Post, error) {
	media.RLock()
	defer media.RUnlock()
	posts := make([]*entity.Post, 0, len(media.posts))

	for _, p := range media.posts {
		if filter == nil || filter.UserID == p.UserID {
			posts = append(posts, toEntityPost(p))
		}
	}

	paginationPosts := SortWithPagination[*entity.Post](posts, pagination)

	return paginationPosts, nil
}

func (media *Media) DeletePost(_ context.Context, ID int) error {
	media.Lock()
	defer media.Unlock()

	if _, ok := media.posts[ID]; !ok {
		return entity.ErrNotFound("post", ID)
	}

	for key := range media.posts[ID].Comments {
		delete(media.posts[ID].Comments, key)
	}

	delete(media.posts, ID)

	return nil
}

func (media *Media) SwitchCommentsState(_ context.Context, ID int, state bool) (*entity.Post, error) {
	media.RLock()
	defer media.RUnlock()
	p, ok := media.posts[ID]
	if !ok {
		return &entity.Post{}, entity.ErrNotFound("post", ID)
	}
	p.IsOpen = state

	return toEntityPost(p), nil
}

func (media *Media) CreateComment(_ context.Context, input entity.Comment) (*entity.Comment, error) {
	media.Lock()
	defer media.Unlock()

	p, ok := media.posts[input.PostID]
	if !ok {
		return &entity.Comment{}, entity.ErrNotFound("post", input.PostID)
	}

	if !p.IsOpen {
		return &entity.Comment{}, entity.ErrPostCommentsDisable
	}

	return p.insertComment(&input)
}

func (media *Media) CreateRepComment(_ context.Context, input entity.Comment) (*entity.Comment, error) {
	media.Lock()
	defer media.Unlock()

	p, ok := media.posts[input.PostID]
	if !ok {
		return &entity.Comment{}, entity.ErrNotFound("post", input.PostID)
	}

	return p.insertRepComment(&input)
}

func (media *Media) GetCommentById(_ context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error) {
	media.RLock()
	defer media.RUnlock()
	comment := &entity.Comment{}

	for _, p := range media.posts {
		if tmp, ok := p.findComment(ID); ok {
			comment = tmp
			break
		}
	}

	if comment == nil {
		return comment, entity.ErrNotFound("comment", ID)
	}

	replies := SortWithPagination[*entity.Comment](comment.Replies, pagination)

	newComment := &entity.Comment{
		ID:        comment.ID,
		ParentID:  comment.ParentID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		Timestamp: comment.Timestamp,
		Replies:   replies,
	}

	return newComment, nil
}

func (media *Media) GetAllComments(_ context.Context, filter *entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error) {
	media.RLock()
	defer media.RUnlock()

	comments := make([]*entity.Comment, 0)

	for _, p := range media.posts {
		comments = append(comments, p.getAllComments()...)
	}

	filterComments := filterComments(comments, filter)

	paginationComments := SortWithPagination[*entity.Comment](filterComments, pagination)
	return paginationComments, nil
}

func filterComments(comments []*entity.Comment, filter *entity.CommentFilter) []*entity.Comment {
	if filter == nil {
		return comments
	}

	var filteredComments []*entity.Comment
	for _, comment := range comments {
		if (filter.PostID == 0 || comment.PostID == filter.PostID) &&
			(filter.UserID == 0 || comment.UserID == filter.UserID) {
			filteredComments = append(filteredComments, comment)
		}
	}

	return filteredComments
}

func (media *Media) DeleteComment(_ context.Context, ID int) error {
	media.Lock()
	defer media.Unlock()

	for _, p := range media.posts {
		if _, ok := p.findComment(ID); ok {
			return p.deleteComment(ID)
		}
	}

	return entity.ErrNotFound("comment", ID)
}
