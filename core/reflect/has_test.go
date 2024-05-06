package reflect

import "testing"

type TestStruct4 struct{}

func (t TestStruct4) MethodX()           {}
func (t TestStruct4) MethodY(a int) bool { return true }

func TestHas(t *testing.T) {
	ts := TestStruct4{}

	tests := []struct {
		name   string
		v      any
		method string
		want   bool
	}{
		{
			name:   "MethodX exists",
			v:      ts,
			method: "MethodX",
			want:   true,
		},
		{
			name:   "MethodY exists",
			v:      ts,
			method: "MethodY",
			want:   true,
		},
		{
			name:   "MethodZ does not exist",
			v:      ts,
			method: "MethodZ",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Has(tt.v, tt.method)
			if got != tt.want {
				t.Errorf("%s: expected %t, got %t", tt.name, tt.want, got)
			}
		})
	}
}
