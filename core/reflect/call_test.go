package reflect

import (
	"fmt"
	"testing"
	"time"

	"github.com/anoideaopen/foundation/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	pb "google.golang.org/protobuf/proto"
)

type TestStruct2 struct{}

func (t *TestStruct2) Method1(ts *time.Time) {
	fmt.Printf("ts: %v\n", ts)
}

func (t *TestStruct2) Method2(ts time.Time) {
	fmt.Printf("ts: %v\n", ts)
}

func (t *TestStruct2) Method3(a *proto.Address) string {
	fmt.Printf("a: %+v\n", a)
	return a.AddrString()
}

func (t *TestStruct2) Method4(in float64) {
	fmt.Printf("in: %+v\n", in)
}

func (t *TestStruct2) Method5(in []float64) {
	fmt.Printf("in: %+v\n", in)
}

func TestCall(t *testing.T) {
	input := &TestStruct2{}

	a := &proto.Address{
		UserID:       "1234",
		Address:      []byte{1, 2, 3, 4},
		IsIndustrial: true,
		IsMultisig:   false,
	}
	aJSON, _ := protojson.Marshal(a)
	aRaw, _ := pb.Marshal(a)

	tests := []struct {
		name      string
		method    string
		arg       string
		wantLen   int
		wantErr   bool
		wantValue any
	}{
		{
			name:    "Method1 with correct time format",
			method:  "Method1",
			arg:     time.Now().Format(time.RFC3339),
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Method2 with correct time format",
			method:  "Method2",
			arg:     time.Now().Format(time.RFC3339),
			wantLen: 0,
			wantErr: false,
		},
		{
			name:      "Method3 with JSON",
			method:    "Method3",
			arg:       string(aJSON),
			wantLen:   1,
			wantErr:   false,
			wantValue: a.AddrString(),
		},
		{
			name:    "Method3 with Protobuf",
			method:  "Method3",
			arg:     string(aRaw),
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "Method4 with float input",
			method:  "Method4",
			arg:     "1234.5678",
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Method5 with array input",
			method:  "Method5",
			arg:     "[1234.5678, 1234.5678]",
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "Method5 with incorrect format",
			method:  "Method5",
			arg:     "1234.5678, 1234.5678",
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Call(input, tt.method, tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Len(t, resp, tt.wantLen)
			if tt.wantValue != nil {
				require.Equal(t, tt.wantValue, resp[0])
			}
		})
	}
}
