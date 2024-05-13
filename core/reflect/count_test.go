package reflect

import "testing"

type TestStruct3 struct{}

func (t TestStruct3) MethodA(x int, y string) bool { return true }
func (t TestStruct3) MethodB() (int, error)        { return 42, nil }
func (t TestStruct3) MethodC(a, b string, c int)   {}

func TestCount(t *testing.T) {
	ts := &TestStruct3{}

	tests := []struct {
		name    string
		v       any
		method  string
		wantIn  int
		wantOut int
	}{
		{
			name:    "MethodA inputs count",
			v:       ts,
			method:  "MethodA",
			wantIn:  2,
			wantOut: 1,
		},
		{
			name:    "MethodB inputs count",
			v:       ts,
			method:  "MethodB",
			wantIn:  0,
			wantOut: 2,
		},
		{
			name:    "MethodC inputs count",
			v:       ts,
			method:  "MethodC",
			wantIn:  3,
			wantOut: 0,
		},
		{
			name:    "MethodUnavailable outputs count",
			v:       ts,
			method:  "MethodUnavailable",
			wantIn:  -1,
			wantOut: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIn, gotOut := Count(tt.v, tt.method)
			if gotIn != tt.wantIn || gotOut != tt.wantOut {
				t.Errorf(
					"%s: gotIn %d, wantIn %d, gotOut %d, wantOut %d",
					tt.name,
					gotIn, tt.wantIn,
					gotOut, tt.wantOut,
				)
			}
		})
	}
}
