package tools

import (
	"GachaServerGin/src"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

func Cache11[K comparable, V any](fn func(K) V) (func(K) V, error) {
	src.Logger.Info("Create Cache11")

	lru := expirable.NewLRU[K, V](64, nil, time.Hour)
	return func(key K) V {
		value, ok := lru.Get(key)
		if !ok {
			value = fn(key)
			lru.Add(key, value)
		}
		return value
	}, nil
}
