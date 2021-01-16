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
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/grpc"
	"solarmesh.io/mockservices/src/service/grpc/protobuf"
)

// server is used to implement protobuf.GreeterServer.
type MockGRpcClient struct {
	AppName string `value:"${app.name}"`

	clientConnector grpc.ClientConnector
}

func newMockGRpcClient(clientConnector grpc.ClientConnector) *MockGRpcClient {
	return &MockGRpcClient{
		clientConnector: clientConnector,
	}
}

func init() {
	app.Register(newMockGRpcClient)
}

// Send implementation
func (s *MockGRpcClient) Send(ctx context.Context, address string, header http.Header) (response *protobuf.MockResponse, err error) {
	//var msc interface{}
	if s.clientConnector == nil {
		err = fmt.Errorf("[src.service.grpc.client] clientConnector is nil")
		log.Error(err)
		return
	}
	conn, e := s.clientConnector.Connect(address)
	if e != nil {
		err = e
		return
	}
	defer conn.Close()

	mockServiceClient := protobuf.NewMockServiceClient(conn)

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	if mockServiceClient != nil {
		// send header to upstream
		md := make(metadata.MD)
		log.Infof("Request Header from GRPC")
		for k, v := range header {
			if strings.Contains(strings.ToLower(k), "fi-") {
				md[k] = v
				log.Infof("> %v: %v", k, v)
				ctx = metadata.AppendToOutgoingContext(ctx, k, v[0])
			}
		}

		// Anything linked to this variable will fetch response headers.
		var responseHeader metadata.MD
		response, err = mockServiceClient.Send(ctx, &protobuf.MockRequest{Host: address}, grpc.Header(&responseHeader))
		log.Infof("Response Header from GRPC")
		for k, v := range responseHeader {
			log.Infof("< %v: %v", k, v)
		}
	}
	return
}
