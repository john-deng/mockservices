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
	"golang.org/x/net/context"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/grpc"
	"solarmesh.io/mockservices/src/rpc/protobuf"
)

// server is used to implement protobuf.GreeterServer.
type MockGRpcServerService struct {

	UpstreamUrls string   `value:"${upstream.urls}"`
	Upstreams    []string `value:"${upstream.urls}"`
	AppName      string   `value:"${app.name}"`
	Version      string   `value:"${app.version}"`
	ClusterName  string   `value:"${cluster.name:my-cluster}"`
	UserData     string   `value:"${user.data:solarmesh}"`

	//mockClientService *client.MockGRpcClientService
}

func newMockServiceServer() protobuf.MockServiceServer {
	return &MockGRpcServerService{}
}

// Send implementation
func (s *MockGRpcServerService) Send(ctx context.Context, request *protobuf.MockRequest) (response *protobuf.MockResponse, err error) {
	response = new(protobuf.MockResponse)

	response.Message = "Success"
	response.Code = 0
	data := new(protobuf.MockData)
	data.Protocol = "GRPC"
	data.App = s.AppName
	data.Version = s.Version
	data.Cluster = s.ClusterName
	data.UserData = s.UserData
	log.Info("send...")
	header, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k, v := range header {
			log.Infof("> %v: %v", k, v)
		}
	}

	response.Data = data

	// Anything linked to this variable will transmit response headers.
	if err := ggrpc.SendHeader(ctx, header); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to send header")
	}

	return
}

func init() {
	// must: register grpc server
	// please note that holaServiceServerImpl must implement protobuf.MockServiceServer, or it won't be registered.
	grpc.Server(protobuf.RegisterMockServiceServer, newMockServiceServer)
}
