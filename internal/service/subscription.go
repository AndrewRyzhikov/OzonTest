package service

import (
	"context"
	"sync"

	"OzonTest/internal/entity"
)

type users map[int]chan *entity.Comment
type posts map[int]users

type SubscriptionService struct {
	sync.RWMutex
	posts posts
}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{posts: make(posts)}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, userID int, postID int) (<-chan *entity.Comment, error) {
	s.Lock()

	if len(s.posts[postID]) == 0 {
		s.posts[postID] = make(users)
	}

	ch := make(chan *entity.Comment)
	s.posts[postID][userID] = ch

	s.Unlock()

	go s.unSubscribe(ctx, ch, userID, postID)

	return ch, nil
}

func (s *SubscriptionService) unSubscribe(ctx context.Context, c chan *entity.Comment, userID int, postID int) {
	<-ctx.Done()

	s.Lock()

	delete(s.posts[postID], userID)
	if len(s.posts[postID]) == 0 {
		delete(s.posts, postID)
	}

	s.Unlock()

	close(c)
}

func (s *SubscriptionService) NotifySubscribers(comment *entity.Comment) {
	s.RLock()
	defer s.RUnlock()

	for _, channel := range s.posts[comment.PostID] {
		//go func() {
		channel <- comment
		//}()
	}
}
