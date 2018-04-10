package message

import (
	"context"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
)

type Request interface {
	UserRequest(string, *pubsub.IdRequest, time.Duration) (*pubsub.UserReply, error)
	UserRequestWithContext(context.Context, string, *pubsub.IdRequest) (*pubsub.UserReply, error)
	IdentityRequest(string, *pubsub.IdentityReq, time.Duration) (*pubsub.IdentityReply, error)
	IdentityRequestWithContext(context.Context, string, *pubsub.IdentityReq) (*pubsub.IdentityReply, error)
}
