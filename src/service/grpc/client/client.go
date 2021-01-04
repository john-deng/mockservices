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

package client

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	grpc2 "hidevops.io/hiboot/pkg/starter/grpc"
	"solarmesh.io/mockservices/src/service/grpc/protobuf"
)

// server is used to implement protobuf.GreeterServer.
type MockGRpcClientService struct {
	AppName string `value:"${app.name}"`

	clientConnector grpc2.ClientConnector
}

func newMockGRpcClientService(clientConnector grpc2.ClientConnector) *MockGRpcClientService {
	return &MockGRpcClientService{
		clientConnector: clientConnector,
	}
}

func init() {
	app.Register(newMockGRpcClientService)
}

// SayMock implements Mockworld.GreeterServer
func (s *MockGRpcClientService) Send(ctx context.Context, address string, header http.Header) (response *protobuf.MockResponse, err error) {
	//var msc interface{}
	var conn *grpc.ClientConn
	conn, err = s.clientConnector.Connect(address)
	if err != nil {
		return
	}
	defer conn.Close()

	log.Infof("gRPC client connected to: %v", address)
	mockServiceClient := protobuf.NewMockServiceClient(conn)

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	if mockServiceClient != nil {
		// send header to upstream
		md := make(metadata.MD)
		for k, v := range header {
			md[k] = v
			log.Infof("> %v: %v", k, v)
			ctx = metadata.AppendToOutgoingContext(ctx, k, v[0])
		}

		// Anything linked to this variable will fetch response headers.
		var responseHeader metadata.MD
		response, err = mockServiceClient.Send(ctx, &protobuf.MockRequest{Host: address}, grpc.Header(&responseHeader))
		for k, v := range responseHeader {
			md[k] = v
			log.Infof("< %v: %v", k, v)
		}
	}
	return
}

