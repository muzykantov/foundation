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

var (
	ErrMethodNotFound          = errors.New("method not found")
	ErrIncorrectArgumentCount  = errors.New("incorrect number of arguments")
	ErrUnsupportedArgumentType = errors.New("unsupported argument type")
)

func Call(input any, methodName string, args ...string) ([]reflect.Value, error) {
	inputVal := reflect.ValueOf(input)

	method := inputVal.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrMethodNotFound, methodName)
	}

	methodType := method.Type()
	if methodType.NumIn() != len(args) {
		return nil, fmt.Errorf("%w: found %d but expected %d", ErrIncorrectArgumentCount, len(args), methodType.NumIn())
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		var (
			argRaw     = []byte(arg)
			argType    = methodType.In(i)
			argPointer = argType.Kind() == reflect.Pointer
		)

		var argValue reflect.Value
		if argPointer {
			argValue = reflect.New(argType.Elem())
			in[i] = argValue
		} else {
			argValue = reflect.New(argType)
			in[i] = argValue.Elem()
		}

		argInterface := argValue.Interface()

		if json.Valid(argRaw) {
			if protoMessage, ok := argInterface.(proto.Message); ok {
				if err := protojson.Unmarshal(argRaw, protoMessage); err != nil {
					return nil, fmt.Errorf("unmarshal JSON for type %s: %w", argType.Name(), err)
				}
			} else {
				if err := json.Unmarshal(argRaw, argInterface); err != nil {
					return nil, fmt.Errorf("unmarshal JSON for type %s: %w", argType.Name(), err)
				}
			}
			continue
		}

		if unmarshaler, ok := argInterface.(encoding.TextUnmarshaler); ok {
			if err := unmarshaler.UnmarshalText(argRaw); err != nil {
				return nil, fmt.Errorf("unmarshal text for type %s: %w", argType.Name(), err)
			}
			continue
		}

		if protoMessage, ok := argInterface.(proto.Message); ok {
			if err := proto.Unmarshal(argRaw, protoMessage); err != nil {
				return nil, fmt.Errorf("unmarshal protobuf for type %s: %w", argType.Name(), err)
			}
			continue
		}

		if unmarshaler, ok := argInterface.(encoding.BinaryUnmarshaler); ok {
			if err := unmarshaler.UnmarshalBinary(argRaw); err != nil {
				return nil, fmt.Errorf("unmarshal binary for type %s: %w", argType.Name(), err)
			}
			continue
		}

		if decoder, ok := argInterface.(gob.GobDecoder); ok {
			if err := decoder.GobDecode(argRaw); err != nil {
				return nil, fmt.Errorf("gob decode for type %s: %w", argType.Name(), err)
			}
			continue
		}

		return nil, fmt.Errorf("%w: %s", ErrUnsupportedArgumentType, argType.Name())
	}

	return method.Call(in), nil
}
