package reflect

import "reflect"

// Count returns the number of input or output parameters of the method specified by the given name.
//
// Parameters:
// - v: the value to check the method on.
// - method: the name of the method.
// - d: the direction to count in (Input or Output).
//
// Returns:
// - int: the number of input or output parameters of the method, or -1 if the method is not found.
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
