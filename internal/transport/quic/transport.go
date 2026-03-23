package quic

import "context"

type Transport struct {
	addr string
}

func New(addr string) *Transport {
	return &Transport{addr: addr}
}

func (t *Transport) Start(context.Context) error {
	return nil
}

func (t *Transport) Stop(context.Context) error {
	return nil
}

func (t *Transport) Addr() string {
	return t.addr
}
