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

type Test struct{}

func (t *Test) Method1(ts *time.Time) {
	fmt.Printf("ts: %v\n", ts)
}

func (t *Test) Method2(ts time.Time) {
	fmt.Printf("ts: %v\n", ts)
}

func (t *Test) Method3(a *proto.Address) string {
	fmt.Printf("a: %+v\n", a)
	return a.AddrString()
}

func (t *Test) Method4(in float64) {
	fmt.Printf("in: %+v\n", in)
}

func (t *Test) Method5(in []float64) {
	fmt.Printf("in: %+v\n", in)
}

func TestCallMethodByName(t *testing.T) {
	input := &Test{}

	resp1, err1 := Call(input, "Method1", time.Now().Format(time.RFC3339))
	require.NoError(t, err1)
	require.Len(t, resp1, 0)

	resp2, err2 := Call(input, "Method2", time.Now().Format(time.RFC3339))
	require.NoError(t, err2)
	require.Len(t, resp2, 0)

	a := &proto.Address{
		UserID:       "1234",
		Address:      []byte{1, 2, 3, 4},
		IsIndustrial: true,
		IsMultisig:   false,
	}
	aJSON, _ := protojson.Marshal(a)

	resp3, err3 := Call(input, "Method3", string(aJSON))
	require.NoError(t, err3)
	require.Len(t, resp3, 1)

	aRaw, _ := pb.Marshal(a)

	resp3, err3 = Call(input, "Method3", string(aRaw))
	require.NoError(t, err3)
	require.Len(t, resp3, 1)

	resp4, err4 := Call(input, "Method4", "1234.5678")
	require.NoError(t, err4)
	require.Len(t, resp4, 0)

	resp5, err5 := Call(input, "Method5", "[1234.5678, 1234.5678]")
	require.NoError(t, err5)
	require.Len(t, resp5, 0)

	resp5, err5 = Call(input, "Method5", "1234.5678, 1234.5678")
	require.Error(t, err5)
	require.Len(t, resp5, 0)
}
