// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// if protoc report command not found error, should install proto and protc-gen-go
// find protoc install instruction on http://google.github.io/proto-lens/installing-protoc.html
// go get -u -v github.com/golang/protobuf/{proto,protoc-gen-go}
//go:generate protoc --proto_path=../protobuf --go_out=plugins=grpc:../protobuf --go_opt=paths=source_relative ../protobuf/mock.proto

package server

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"solarmesh.io/mockservices/src/model"
	"solarmesh.io/mockservices/src/service"
	"solarmesh.io/mockservices/src/service/grpc/protobuf"
)

// server is used to implement protobuf.GreeterServer.
type MockGRpcServerService struct {
	mockService *service.MockService
}

func newMockServiceServer(mockService *service.MockService) protobuf.MockServiceServer {
	return &MockGRpcServerService{mockService: mockService}
}

func init() {
	// must: register grpc server
	// please note that holaServiceServerImpl must implement protobuf.MockServiceServer, or it won't be registered.
	grpc.Server(protobuf.RegisterMockServiceServer, newMockServiceServer)
}


// Send implementation
func (s *MockGRpcServerService) Send(ctx context.Context, request *protobuf.MockRequest) (response *protobuf.MockResponse, err error) {
	httpHeader := make(http.Header)
	header, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k, v := range header {
			log.Infof("> %v: %v", k, v)
			httpHeader.Set(k, v[0])
		}
	}
	var resp *model.Response
	response = new(protobuf.MockResponse)
	resp, err = s.mockService.SendRequest("GRPC", nil, httpHeader)
	if err == nil {
		b, e := json.Marshal(resp)
		if e == nil {
			e = json.Unmarshal(b, response)
		}
	}

	// Anything linked to this variable will transmit response headers.
	if err := ggrpc.SendHeader(ctx, header); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to send header")
	}

	return
}
