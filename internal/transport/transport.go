package transport

import "context"

type Kind string

const (
	KindTCP  Kind = "tcp"
	KindQUIC Kind = "quic"
)

type Transport interface {
	Start(context.Context) error
	Stop(context.Context) error
	Addr() string
}
