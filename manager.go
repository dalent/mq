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
	Consume(act Action, ack bool) error
	Publish(string, []byte) error
	NewWork(string, []byte) error
	Worker(int, Action) error
}

const (
	RABBIT = "rabbit"
)

type New func() Provider

//keep the provider and the interface
var mq = make(map[string]New)

func Register(name string, function New) {
	if _, ok := mq[name]; ok {
		panic("name already registered:" + name)
	}
	mq[name] = function
}

func NewProvider(name string) (Provider, error) {
	if v, ok := mq[name]; ok {
		return v(), nil
	}

	return nil, errors.New("not support manager")
}
