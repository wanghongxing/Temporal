// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"sync"

	"github.com/RTradeLtd/grpc/pay"
	"github.com/RTradeLtd/grpc/pay/request"
	"github.com/RTradeLtd/grpc/pay/response"
	"google.golang.org/grpc"
)

type FakeSignerClient struct {
	GetSignedMessageStub        func(context.Context, *request.SignRequest, ...grpc.CallOption) (*response.SignResponse, error)
	getSignedMessageMutex       sync.RWMutex
	getSignedMessageArgsForCall []struct {
		arg1 context.Context
		arg2 *request.SignRequest
		arg3 []grpc.CallOption
	}
	getSignedMessageReturns struct {
		result1 *response.SignResponse
		result2 error
	}
	getSignedMessageReturnsOnCall map[int]struct {
		result1 *response.SignResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSignerClient) GetSignedMessage(arg1 context.Context, arg2 *request.SignRequest, arg3 ...grpc.CallOption) (*response.SignResponse, error) {
	fake.getSignedMessageMutex.Lock()
	ret, specificReturn := fake.getSignedMessageReturnsOnCall[len(fake.getSignedMessageArgsForCall)]
	fake.getSignedMessageArgsForCall = append(fake.getSignedMessageArgsForCall, struct {
		arg1 context.Context
		arg2 *request.SignRequest
		arg3 []grpc.CallOption
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetSignedMessage", []interface{}{arg1, arg2, arg3})
	fake.getSignedMessageMutex.Unlock()
	if fake.GetSignedMessageStub != nil {
		return fake.GetSignedMessageStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getSignedMessageReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSignerClient) GetSignedMessageCallCount() int {
	fake.getSignedMessageMutex.RLock()
	defer fake.getSignedMessageMutex.RUnlock()
	return len(fake.getSignedMessageArgsForCall)
}

func (fake *FakeSignerClient) GetSignedMessageCalls(stub func(context.Context, *request.SignRequest, ...grpc.CallOption) (*response.SignResponse, error)) {
	fake.getSignedMessageMutex.Lock()
	defer fake.getSignedMessageMutex.Unlock()
	fake.GetSignedMessageStub = stub
}

func (fake *FakeSignerClient) GetSignedMessageArgsForCall(i int) (context.Context, *request.SignRequest, []grpc.CallOption) {
	fake.getSignedMessageMutex.RLock()
	defer fake.getSignedMessageMutex.RUnlock()
	argsForCall := fake.getSignedMessageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeSignerClient) GetSignedMessageReturns(result1 *response.SignResponse, result2 error) {
	fake.getSignedMessageMutex.Lock()
	defer fake.getSignedMessageMutex.Unlock()
	fake.GetSignedMessageStub = nil
	fake.getSignedMessageReturns = struct {
		result1 *response.SignResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSignerClient) GetSignedMessageReturnsOnCall(i int, result1 *response.SignResponse, result2 error) {
	fake.getSignedMessageMutex.Lock()
	defer fake.getSignedMessageMutex.Unlock()
	fake.GetSignedMessageStub = nil
	if fake.getSignedMessageReturnsOnCall == nil {
		fake.getSignedMessageReturnsOnCall = make(map[int]struct {
			result1 *response.SignResponse
			result2 error
		})
	}
	fake.getSignedMessageReturnsOnCall[i] = struct {
		result1 *response.SignResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeSignerClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getSignedMessageMutex.RLock()
	defer fake.getSignedMessageMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSignerClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ pay.SignerClient = new(FakeSignerClient)
