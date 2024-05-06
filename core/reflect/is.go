package reflect

import (
	"reflect"
)

// Is checks if the target type matches the type of the n-th argument or return value of the method in the given direction.
//
// Parameters:
// - v: the value to check the method on.
// - method: the name of the method.
// - n: the index of the argument or return value to check.
// - d: the direction to check in (Input or Output).
// - target: the target type to compare against.
//
// Returns:
// - bool: true if the target type matches the type of the argument or return value, false otherwise.
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
