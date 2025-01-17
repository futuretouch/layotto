/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wasm

import (
	"context"
	anypb "github.com/golang/protobuf/ptypes/any"
	"mosn.io/layotto/pkg/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm010"
	"mosn.io/proxy-wasm-go-host/proxywasm/common"
	proxywasm "mosn.io/proxy-wasm-go-host/proxywasm/v1"
)

// LayottoHandler implement proxywasm.ImportsHandler
type LayottoHandler struct {
	proxywasm010.DefaultImportsHandler

	IoBuffer common.IoBuffer
}

var _ proxywasm.ImportsHandler = &LayottoHandler{}

var Layotto grpc.API

func (d *LayottoHandler) GetState(storeName string, key string) (string, proxywasm.WasmResult) {
	req := &runtimev1pb.GetStateRequest{
		StoreName: storeName,
		Key:       key,
	}
	resp, err := Layotto.GetState(context.Background(), req)
	if err != nil {
		return "", proxywasm.WasmResultInternalFailure
	}
	return string(resp.Data), proxywasm.WasmResultOk
}

func (d *LayottoHandler) InvokeService(id string, method string, param string) (string, proxywasm.WasmResult) {
	req := &runtimev1pb.InvokeServiceRequest{
		Id: id,
		Message: &runtimev1pb.CommonInvokeRequest{
			Method: method,
			Data:   &anypb.Any{Value: []byte(param)},
		},
	}
	resp, err := Layotto.InvokeService(context.Background(), req)
	if err != nil {
		return "", proxywasm.WasmResultInternalFailure
	}
	return string(resp.Data.Value), proxywasm.WasmResultOk
}

func (d *LayottoHandler) GetFuncCallData() common.IoBuffer {
	if d.IoBuffer == nil {
		d.IoBuffer = common.NewIoBufferBytes(make([]byte, 0))
	}
	return d.IoBuffer
}
