package contracts

import (
	"context"

	"OzonTest/internal/entity"
)

type Subscription interface {
	Subscribe(ctx context.Context, userID int, postID int) (<-chan *entity.Comment, error)
	NotifySubscribers(comment *entity.Comment)
}
