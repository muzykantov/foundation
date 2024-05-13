package reflect

import "reflect"

// Count counts the number of input and output parameters of a method in a given value.
//
// Parameters:
// - v: The value in which the method is located.
// - method: The name of the method.
//
// Returns:
// - in: The number of input parameters of the method.
// - out: The number of output parameters of the method.
func Count(v any, method string) (in int, out int) {
	methodVal := reflect.ValueOf(v).MethodByName(method)
	if !methodVal.IsValid() {
		return -1, -1
	}

	methodType := methodVal.Type()
	return methodType.NumIn(), methodType.NumOut()
}
