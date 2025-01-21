package typed

import "github.com/kmou424/ero"

type Map[K comparable] map[K]Value

func NewMap[K comparable]() Map[K] {
	return make(Map[K])
}

func (tm Map[K]) Set(key K, value any) {
	tm[key] = newTypedValue(value)
}

func (tm Map[K]) Get(key K) (value any, ok bool) {
	tv, ok := tm[key]
	if ok {
		value = tv.v
	}
	return value, ok
}

func (tm Map[K]) MustGet(key K) any {
	value, ok := tm.Get(key)
	if !ok {
		panic(ero.Newf("key %v not found in map", key))
	}
	return value
}
