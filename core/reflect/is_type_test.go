package reflect

import (
	"errors"
	"testing"
)

type TestStruct1 struct{}

func (t *TestStruct1) MethodWithPointer(arg *TestStruct1) *TestStruct1 { return arg }
func (t *TestStruct1) MethodWithSlice(arg []int) []int                 { return arg }
func (t *TestStruct1) MethodWithValue(arg int) int                     { return arg }
func (t *TestStruct1) MethodWithError() error                          { return errors.New("test error") }
func (t *TestStruct1) MethodWithBool() bool                            { return true }

func TestIsType(t *testing.T) {
	testStruct := &TestStruct1{}

	tests := []struct {
		name        string
		v           any
		method      string
		n           int
		target      any
		checkReturn bool
		want        bool
	}{
		{
			name:        "Type of argument for MethodWithPointer",
			v:           testStruct,
			method:      "MethodWithPointer",
			n:           0,
			target:      &TestStruct1{},
			checkReturn: false,
			want:        true,
		},
		{
			name:        "Type of argument for MethodWithPointer",
			v:           testStruct,
			method:      "MethodWithPointer",
			n:           0,
			target:      (*TestStruct1)(nil),
			checkReturn: false,
			want:        true,
		},
		{
			name:        "Type of return value for MethodWithPointer",
			v:           testStruct,
			method:      "MethodWithPointer",
			n:           0,
			target:      &TestStruct1{},
			checkReturn: true,
			want:        true,
		},
		{
			name:        "Type of return value for MethodWithError",
			v:           testStruct,
			method:      "MethodWithError",
			n:           0,
			target:      errors.New(""),
			checkReturn: true,
			want:        true,
		},
		{
			name:        "Type of return value for MethodWithBool",
			v:           testStruct,
			method:      "MethodWithBool",
			n:           0,
			target:      true,
			checkReturn: true,
			want:        true,
		},
		{
			name:        "Type of argument for MethodWithSlice",
			v:           testStruct,
			method:      "MethodWithSlice",
			n:           0,
			target:      []int{},
			checkReturn: false,
			want:        true,
		},
		{
			name:        "Type of argument for MethodWithSlice",
			v:           testStruct,
			method:      "MethodWithSlice",
			n:           0,
			target:      nil,
			checkReturn: false,
			want:        false,
		},
		{
			name:        "Type of return value for MethodWithSlice",
			v:           testStruct,
			method:      "MethodWithSlice",
			n:           0,
			target:      []int{},
			checkReturn: true,
			want:        true,
		},
		{
			name:        "Type of argument for MethodWithValue",
			v:           testStruct,
			method:      "MethodWithValue",
			n:           0,
			target:      0, // int
			checkReturn: false,
			want:        true,
		},
		{
			name:        "Type of return value for MethodWithValue",
			v:           testStruct,
			method:      "MethodWithValue",
			n:           0,
			target:      0, // int
			checkReturn: true,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsType(tt.v, tt.method, tt.n, tt.checkReturn, tt.target); got != tt.want {
				t.Errorf("IsType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
