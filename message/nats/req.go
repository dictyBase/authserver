package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/dictyBase/authserver/message"
	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"
	gnats "github.com/nats-io/go-nats"

	"github.com/nats-io/go-nats/encoders/protobuf"
)

type natsRequest struct {
	econn *gnats.EncodedConn
}

func NewRequest(host, port string, options ...gnats.Option) (message.Request, error) {
	nc, err := gnats.Connect(fmt.Sprintf("nats://%s:%s", host, port), options...)
	if err != nil {
		return &natsRequest{}, err
	}
	ec, err := gnats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return &natsRequest{}, err
	}
	return &natsRequest{econn: ec}, nil
}

func (n *natsRequest) UserRequest(subj string, r *pubsub.IdRequest, timeout time.Duration) (*pubsub.UserReply, error) {
	reply := &pubsub.UserReply{}
	err := n.econn.Request(subj, r, reply, timeout)
	return reply, err
}

func (n *natsRequest) UserRequestWithContext(ctx context.Context, subj string, r *pubsub.IdRequest) (*pubsub.UserReply, error) {
	reply := &pubsub.UserReply{}
	err := n.econn.RequestWithContext(ctx, subj, r, reply)
	return reply, err
}

func (n *natsRequest) IdentityRequest(subj string, r *pubsub.IdentityReq, timeout time.Duration) (*pubsub.IdentityReply, error) {
	reply := &pubsub.IdentityReply{}
	err := n.econn.Request(subj, r, reply, timeout)
	return reply, err
}

func (n *natsRequest) IdentityRequestWithContext(ctx context.Context, subj string, r *pubsub.IdentityReq) (*pubsub.IdentityReply, error) {
	reply := &pubsub.IdentityReply{}
	err := n.econn.RequestWithContext(ctx, subj, r, reply)
	return reply, err
}

func (n *natsRequest) IsActive() bool {
	return n.econn.Conn.IsConnected()
}
