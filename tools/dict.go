package tools

type DefaultDict[K comparable, V any] struct {
	data         map[K]V
	defaultValue func() V
}

func (dict DefaultDict[K, V]) Get(key K) V {
	value, has := dict.data[key]
	if !has {
		value = dict.defaultValue()
		dict.data[key] = value
	}
	return value
}

func (dict DefaultDict[K, V]) Set(key K, value V) {
	dict.data[key] = value
}

func (dict DefaultDict[K, V]) Data() map[K]V {
	return dict.data
}

func NewDefaultDict[K comparable, V any](f func() V) DefaultDict[K, V] {
	return DefaultDict[K, V]{
		data:         map[K]V{},
		defaultValue: f,
	}
}

type Counter[K comparable] struct {
	DefaultDict[K, int]
}

func (c Counter[K]) Inc(key K, value int) {
	c.Set(key, c.Get(key)+value)
}

func NewCounter[K comparable]() Counter[K] {
	return Counter[K]{
		NewDefaultDict[K](func() int {
			return 0
		}),
	}
}
