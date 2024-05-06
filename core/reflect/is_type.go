package reflect

import (
	"reflect"
)

func IsType(v any, method string, offset int, checkReturn bool, target any) bool {
	methodVal := reflect.ValueOf(v).MethodByName(method)

	if !methodVal.IsValid() {
		return false
	}

	var (
		methodType = methodVal.Type()
		targetType = reflect.TypeOf(target)
	)
	if checkReturn {
		if offset >= methodType.NumOut() {
			return false
		}
		outType := methodType.Out(offset)
		return targetType == outType || (targetType != nil && targetType.Implements(outType))
	} else {
		if offset >= methodType.NumIn() {
			return false
		}
		inType := methodType.In(offset)
		return targetType == inType || (targetType != nil && targetType.Implements(inType))
	}
}
