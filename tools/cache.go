package tools

import (
	"GachaServerGin/src"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

func Cache11[K comparable, V any](fn func(K) V) func(K) V {
	src.Logger.Info("Create Cache11")

	lru := expirable.NewLRU[K, V](64, nil, time.Hour)
	return func(key K) V {
		value, ok := lru.Get(key)
		if !ok {
			value = fn(key)
			lru.Add(key, value)
		}
		return value
	}
}

type Cop[A any, B any] struct {
	va A
	vb B
}

func Cache12[K comparable, A any, B any](fn func(K) (A, B)) func(K) (A, B) {
	src.Logger.Info("Create Cache12")
	c11 := Cache11(func(key K) Cop[A, B] {
		a, b := fn(key)
		return Cop[A, B]{
			va: a,
			vb: b,
		}
	})
	return func(key K) (A, B) {
		cop := c11(key)
		return cop.va, cop.vb
	}
}

func Cache11e[K comparable, V any](fn func(K) (V, error)) func(K) (V, error) {
	src.Logger.Info("Create Cache11e")
	lru := expirable.NewLRU[K, V](64, nil, time.Hour)
	return func(key K) (V, error) {
		var (
			value V
			err   error
			ok    bool
		)
		value, ok = lru.Get(key)
		if !ok {
			value, err = fn(key)
			if err != nil {
				return value, err
			}
			lru.Add(key, value)
		}
		return value, nil
	}
}
