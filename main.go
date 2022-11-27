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

package main

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"github.com/hidevopsio/hiboot/pkg/starter/httpclient"
	"github.com/hidevopsio/hiboot/pkg/starter/jaeger"
	"github.com/hidevopsio/hiboot/pkg/starter/locale"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	_ "github.com/john-deng/mockservices/src/controller"
	_ "github.com/john-deng/mockservices/src/service/grpc/server"
	"github.com/john-deng/mockservices/src/service/tcp"
)

func main() {
	// create new web application and run it
	app.IncludeProfiles(
		logging.Profile,
		locale.Profile,
		grpc.Profile,
		tcp.Profile,
		jaeger.Profile,
		httpclient.Profile,
		actuator.Profile)

	web.NewApplication().
		Run()
}
