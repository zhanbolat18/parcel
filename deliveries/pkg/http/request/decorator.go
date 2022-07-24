package request

import (
	"context"
	"errors"
	"net/http"
)

type RequestDecorator interface {
	Decorate(req *http.Request)
}

type DecoratorMap interface {
	RequestDecorator
	Store(key interface{}, value RequestDecorator) error
	Delete(key interface{})
}

type headerProxy struct {
	dec        RequestDecorator
	key, value string
}

func (t *headerProxy) Decorate(req *http.Request) {
	t.dec.Decorate(req)
	req.Header.Add(t.key, t.value)
}

type emptyDecorator struct {
}

func (e *emptyDecorator) Decorate(req *http.Request) {
}

type decoratorMaps struct {
	m map[context.Context]RequestDecorator
}

func (d *decoratorMaps) Store(key interface{}, value RequestDecorator) error {
	if ctx, ok := key.(context.Context); ok {
		d.m[ctx] = value
		return nil
	}
	return errors.New("key mus implement Context interface")
}

func (d *decoratorMaps) Delete(key interface{}) {
	delete(d.m, key.(context.Context))
}

func (d *decoratorMaps) Decorate(req *http.Request) {
	if dec, ok := d.m[req.Context()]; ok {
		dec.Decorate(req)
	}
}

func NewDecoratorMaps() DecoratorMap {
	return &decoratorMaps{m: map[context.Context]RequestDecorator{}}
}

func ApiAuthProxyDecorator(dec RequestDecorator, fullToken string) RequestDecorator {
	return &headerProxy{key: "Authorization", value: fullToken, dec: dec}
}

func EmptyDecorator() RequestDecorator {
	return &emptyDecorator{}
}
