package transport

import "context"

type Kind string

const (
	KindTCP  Kind = "tcp"
	KindQUIC Kind = "quic"
)

type MessageType string

const (
	MessageReplicateSet MessageType = "replicate_set"
	MessageFetchValue   MessageType = "fetch_value"
	MessagePing         MessageType = "ping"
)

type Request struct {
	Type  MessageType `json:"type"`
	From  string      `json:"from"`
	Key   string      `json:"key,omitempty"`
	Value string      `json:"value,omitempty"`
}

type Response struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	From    string `json:"from,omitempty"`
	Value   string `json:"value,omitempty"`
	Found   bool   `json:"found,omitempty"`
}

type HandlerFunc func(context.Context, Request) (Response, error)

type Transport interface {
	Start(context.Context, HandlerFunc) error
	Stop(context.Context) error
	Send(context.Context, string, Request) (Response, error)
	Addr() string
	Kind() Kind
}
