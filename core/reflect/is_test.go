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

func TestIs(t *testing.T) {
	testStruct := &TestStruct1{}

	tests := []struct {
		name   string
		v      any
		method string
		n      int
		output bool
		target any
		want   bool
	}{
		{
			name:   "Type of argument for MethodWithPointer",
			v:      testStruct,
			method: "MethodWithPointer",
			n:      0,
			output: false,
			target: &TestStruct1{},
			want:   true,
		},
		{
			name:   "Type of argument for MethodWithPointer",
			v:      testStruct,
			method: "MethodWithPointer",
			n:      0,
			output: false,
			target: (*TestStruct1)(nil),
			want:   true,
		},
		{
			name:   "Type of return value for MethodWithPointer",
			v:      testStruct,
			method: "MethodWithPointer",
			n:      0,
			output: true,
			target: &TestStruct1{},
			want:   true,
		},
		{
			name:   "Type of return value for MethodWithError",
			v:      testStruct,
			method: "MethodWithError",
			n:      0,
			output: true,
			target: errors.New(""),
			want:   true,
		},
		{
			name:   "Type of return value for MethodWithBool",
			v:      testStruct,
			method: "MethodWithBool",
			n:      0,
			output: true,
			target: true,
			want:   true,
		},
		{
			name:   "Type of argument for MethodWithSlice",
			v:      testStruct,
			method: "MethodWithSlice",
			n:      0,
			output: false,
			target: []int{},
			want:   true,
		},
		{
			name:   "Type of argument for MethodWithSlice",
			v:      testStruct,
			method: "MethodWithSlice",
			n:      0,
			output: false,
			target: nil,
			want:   false,
		},
		{
			name:   "Type of return value for MethodWithSlice",
			v:      testStruct,
			method: "MethodWithSlice",
			n:      0,
			output: true,
			target: []int{},
			want:   true,
		},
		{
			name:   "Type of argument for MethodWithValue",
			v:      testStruct,
			method: "MethodWithValue",
			n:      0,
			output: false,
			target: 0, // int
			want:   true,
		},
		{
			name:   "Type of return value for MethodWithValue",
			v:      testStruct,
			method: "MethodWithValue",
			n:      0,
			output: true,
			target: 0, // int
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.v, tt.method, tt.n, tt.output, tt.target); got != tt.want {
				t.Errorf("IsType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
