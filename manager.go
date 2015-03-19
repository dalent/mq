package mq

import (
	"errors"
)

type Action interface {
	Do(bytes []byte)
}

type Provider interface {
	SetUp(url, queue string) error
	Close()
	Consume(act Action) error
	Publish(string, []byte) error
	NewWork(string, []byte) error
	Worker(int, Action) error
}

const (
	RABBIT = "rabbit"
)

//keep the provider and the interface
var mq = make(map[string]Provider)

func Register(name string, provider Provider) {
	mq[name] = provider
}

func NewProvider(name string) (provider Provider, e error) {
	if provider, ok := mq[name]; ok {
		return provider, nil
	}

	return nil, errors.New("not support manager")
}
