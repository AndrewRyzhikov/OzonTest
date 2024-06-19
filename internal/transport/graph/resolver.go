package graph

import (
	"OzonTest/internal/service/contracts"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CommentService      contracts.Comment
	PostService         contracts.Post
	SubscriptionService contracts.Subscription
}

func NewResolver(commentService contracts.Comment, postService contracts.Post, subscriptionService contracts.Subscription) *Resolver {
	return &Resolver{CommentService: commentService, PostService: postService, SubscriptionService: subscriptionService}
}
