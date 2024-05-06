package reflect

import "reflect"

// Has checks if the given value has a method with the specified name.
//
// Parameters:
// - v: the value to check for the method.
// - method: the name of the method to check for.
//
// Returns:
// - bool: true if the value has the method, false otherwise.
func Has(v any, method string) bool {
	return reflect.ValueOf(v).MethodByName(method).IsValid()
}
