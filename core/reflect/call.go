package reflect

import (
	"encoding"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Error types.
var (
	ErrMethodNotFound          = errors.New("method not found")
	ErrIncorrectArgumentCount  = errors.New("incorrect number of arguments")
	ErrUnsupportedArgumentType = errors.New("unsupported argument type")
	ErrInvalidArgumentValue    = errors.New("invalid argument value")
)

func Call(v any, method string, args ...string) ([]any, error) {
	inputVal := reflect.ValueOf(v)

	methodVal := inputVal.MethodByName(method)
	if !methodVal.IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrMethodNotFound, method)
	}

	methodType := methodVal.Type()
	if methodType.NumIn() != len(args) {
		return nil, fmt.Errorf("%w: found %d but expected %d", ErrIncorrectArgumentCount, len(args), methodType.NumIn())
	}

	var (
		in  = make([]reflect.Value, len(args))
		err error
	)
	for i, arg := range args {
		if in[i], err = value(arg, methodType.In(i)); err != nil {
			return nil, fmt.Errorf("%w: call %s", err, method)
		}
	}

	output := make([]any, methodType.NumOut())
	for i, res := range methodVal.Call(in) {
		output[i] = res.Interface()
	}

	return output, nil
}

func value(s string, t reflect.Type) (reflect.Value, error) {
	var (
		argRaw     = []byte(s)
		argType    = t
		argPointer = t.Kind() == reflect.Pointer
	)

	// if the type is a pointer, create a new instance of the element type
	// and assign it to the pointer. otherwise, create a new instance of the
	// type and assign it to the element of the newly created pointer.
	// this is necessary because we can't call json.Unmarshal directly on a
	// pointer.
	var (
		argValue reflect.Value
		outValue reflect.Value
	)
	if argPointer {
		argValue = reflect.New(t.Elem()) // create a new instance of the element type
		outValue = argValue              // assign the pointer to the return value
	} else {
		argValue = reflect.New(t)  // create a new instance of the type
		outValue = argValue.Elem() // assign the element of the pointer to the return value
	}

	argInterface := argValue.Interface()

	if json.Valid(argRaw) {
		var err error
		if protoMessage, ok := argInterface.(proto.Message); ok {
			err = protojson.Unmarshal(argRaw, protoMessage)
		} else {
			err = json.Unmarshal(argRaw, argInterface)
		}
		if err != nil {
			return outValue, fmt.Errorf("%w: unmarshal JSON for type %s: %v", ErrInvalidArgumentValue, argType.Name(), err)
		}

		return outValue, nil
	}

	if unmarshaler, ok := argInterface.(encoding.TextUnmarshaler); ok {
		if err := unmarshaler.UnmarshalText(argRaw); err != nil {
			return outValue, fmt.Errorf("%w: unmarshal text for type %s: %v", ErrInvalidArgumentValue, argType.Name(), err)
		}

		return outValue, nil
	}

	if protoMessage, ok := argInterface.(proto.Message); ok {
		if err := proto.Unmarshal(argRaw, protoMessage); err != nil {
			return outValue, fmt.Errorf("%w: unmarshal protobuf for type %s: %v", ErrInvalidArgumentValue, argType.Name(), err)
		}

		return outValue, nil
	}

	if unmarshaler, ok := argInterface.(encoding.BinaryUnmarshaler); ok {
		if err := unmarshaler.UnmarshalBinary(argRaw); err != nil {
			return outValue, fmt.Errorf("%w: unmarshal binary for type %s: %v", ErrInvalidArgumentValue, argType.Name(), err)
		}

		return outValue, nil
	}

	if decoder, ok := argInterface.(gob.GobDecoder); ok {
		if err := decoder.GobDecode(argRaw); err != nil {
			return outValue, fmt.Errorf("%w: gob decode for type %s: %v", ErrInvalidArgumentValue, argType.Name(), err)
		}

		return outValue, nil
	}

	return outValue, fmt.Errorf("%w: type %s", ErrUnsupportedArgumentType, argType.Name())
}
