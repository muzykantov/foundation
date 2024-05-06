package reflect

import (
	"reflect"
)

func Is(v any, method string, n int, d Direction, target any) bool {
	methodVal := reflect.ValueOf(v).MethodByName(method)

	if !methodVal.IsValid() {
		return false
	}

	var (
		methodType = methodVal.Type()
		targetType = reflect.TypeOf(target)
	)

	switch d {
	case Input:
		if n >= methodType.NumIn() {
			return false
		}

		inType := methodType.In(n)
		return targetType == inType || (targetType != nil && targetType.Implements(inType))

	case Output:
		if n >= methodType.NumOut() {
			return false
		}

		outType := methodType.Out(n)
		return targetType == outType || (targetType != nil && targetType.Implements(outType))

	default:
		return false
	}
}
