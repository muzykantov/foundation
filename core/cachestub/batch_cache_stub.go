package cachestub

import (
	"net/http"
	"strings"

	"github.com/anoideaopen/foundation/proto"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type BatchCacheStub struct {
	shim.ChaincodeStubInterface
	batchWriteCache   map[string]*proto.WriteElement
	batchReadeCache   map[string]*proto.WriteElement
	invokeResultCache map[string]peer.Response
	Swaps             []*proto.Swap
	MultiSwaps        []*proto.MultiSwap
}

func NewBatchCacheStub(stub shim.ChaincodeStubInterface) *BatchCacheStub {
	return &BatchCacheStub{
		ChaincodeStubInterface: stub,
		batchWriteCache:        make(map[string]*proto.WriteElement),
		batchReadeCache:        make(map[string]*proto.WriteElement),
		invokeResultCache:      make(map[string]peer.Response),
	}
}

// GetState returns state from BatchCacheStub cache or, if absent, from chaincode state
func (bs *BatchCacheStub) GetState(key string) ([]byte, error) {
	if existsElement, ok := bs.batchWriteCache[key]; ok {
		return existsElement.GetValue(), nil
	}

	if existsElement, ok := bs.batchReadeCache[key]; ok {
		return existsElement.GetValue(), nil
	}

	value, err := bs.ChaincodeStubInterface.GetState(key)
	if err != nil {
		return nil, err
	}

	bs.batchReadeCache[key] = &proto.WriteElement{Key: key, Value: value}

	return value, nil
}

// PutState puts state to a BatchCacheStub cache
func (bs *BatchCacheStub) PutState(key string, value []byte) error {
	bs.batchWriteCache[key] = &proto.WriteElement{Key: key, Value: value}
	return nil
}

// Commit puts state from a BatchCacheStub cache to the chaincode state
func (bs *BatchCacheStub) Commit() error {
	for key, element := range bs.batchWriteCache {
		if element.GetIsDeleted() {
			if err := bs.ChaincodeStubInterface.DelState(key); err != nil {
				return err
			}
		} else {
			if err := bs.ChaincodeStubInterface.PutState(key, element.GetValue()); err != nil {
				return err
			}
		}
	}
	return nil
}

// DelState - marks state in BatchCacheStub cache as deleted
func (bs *BatchCacheStub) DelState(key string) error {
	bs.batchWriteCache[key] = &proto.WriteElement{Key: key, IsDeleted: true}
	return nil
}

func (bs *BatchCacheStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) peer.Response {
	keys := []string{channel, chaincodeName}
	for _, arg := range args {
		keys = append(keys, string(arg))
	}
	key := strings.Join(keys, "")

	if result, ok := bs.invokeResultCache[key]; ok {
		return result
	}

	resp := bs.ChaincodeStubInterface.InvokeChaincode(chaincodeName, args, channel)

	if resp.GetStatus() == http.StatusOK && len(resp.GetPayload()) != 0 { //nolint:nestif
		if string(args[0]) == "checkKeys" {
			func() {
				addrMsg := &proto.AclResponse{}
				if err := pb.Unmarshal(resp.GetPayload(), addrMsg); err != nil {
					return
				}

				addr := addrMsg.GetAddress().GetAddress().AddrString()

				address, err := pb.Marshal(addrMsg.GetAddress().GetAddress())
				if err != nil {
					return
				}
				keyCheck := channel + chaincodeName + "checkAddress" + addr
				bs.invokeResultCache[keyCheck] = peer.Response{
					Status:  http.StatusOK,
					Payload: address,
				}

				accinfo, err := pb.Marshal(addrMsg.GetAccount())
				if err != nil {
					return
				}
				keyAccInfo := channel + chaincodeName + "getAccountInfo" + addr
				bs.invokeResultCache[keyAccInfo] = peer.Response{
					Status:  http.StatusOK,
					Payload: accinfo,
				}
			}()
		}

		bs.invokeResultCache[key] = resp
	}

	return resp
}
