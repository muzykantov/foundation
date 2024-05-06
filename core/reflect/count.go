package reflect

import "reflect"

func Count(v any, method string, d Direction) int {
	methodVal := reflect.ValueOf(v).MethodByName(method)
	if !methodVal.IsValid() {
		return -1
	}

	methodType := methodVal.Type()

	switch d {
	case Input:
		return methodType.NumIn()
	case Output:
		return methodType.NumOut()
	default:
		return -1
	}
}
