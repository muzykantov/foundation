package reflect

import "testing"

type TestStruct3 struct{}

func (t TestStruct3) MethodA(x int, y string) bool { return true }
func (t TestStruct3) MethodB() (int, error)        { return 42, nil }
func (t TestStruct3) MethodC(a, b string, c int)   {}

func TestCount(t *testing.T) {
	ts := TestStruct3{}

	tests := []struct {
		name   string
		v      any
		method string
		d      Direction
		want   int
	}{
		{
			name:   "MethodA inputs count",
			v:      ts,
			method: "MethodA",
			d:      Input,
			want:   2, // x и y
		},
		{
			name:   "MethodA outputs count",
			v:      ts,
			method: "MethodA",
			d:      Output,
			want:   1, // bool
		},
		{
			name:   "MethodB inputs count",
			v:      ts,
			method: "MethodB",
			d:      Input,
			want:   0,
		},
		{
			name:   "MethodB outputs count",
			v:      ts,
			method: "MethodB",
			d:      Output,
			want:   2, // int и error
		},
		{
			name:   "MethodC inputs count",
			v:      ts,
			method: "MethodC",
			d:      Input,
			want:   3, // a, b, и c
		},
		{
			name:   "MethodC outputs count",
			v:      ts,
			method: "MethodC",
			d:      Output,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Count(tt.v, tt.method, tt.d)
			if got != tt.want {
				t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}
