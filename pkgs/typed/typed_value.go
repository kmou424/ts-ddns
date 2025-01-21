package typed

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/kmou424/ero"
	"reflect"
)

type Value struct {
	v any
}

func newTypedValue(v any) Value {
	return Value{v: v}
}

func (tv Value) ToInt() int {
	ret, err := mathutil.ToInt(tv.v)
	if err != nil {
		return 0
	}
	return ret
}

func (tv Value) ToFloat() float64 {
	ret, err := mathutil.ToFloat(tv.v)
	if err != nil {
		return 0
	}
	return ret
}

func (tv Value) ToBool() bool {
	ret, err := goutil.ToBool(tv.v)
	if err != nil {
		return false
	}
	return ret
}

func (tv Value) ToString() string {
	ret, err := mathutil.ToString(tv.v)
	if err != nil {
		return ""
	}
	return ret
}

func ValueTo[T any](v Value) (casted T, err error) {
	var ok bool
	if casted, ok = v.v.(T); !ok {
		refType := reflect.TypeOf(casted)
		err = ero.Newf("cannot cast value to type %s", refType)
		return
	}
	return
}
