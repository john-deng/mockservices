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
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	_ "hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/jaeger"
	"hidevops.io/hiboot/pkg/starter/logging"
	_ "solarmesh.io/mockservices/src/controller"
	_ "solarmesh.io/mockservices/src/service/grpc/server"
	_ "solarmesh.io/mockservices/src/service/tcp"
)

func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude, logging.Profile, jaeger.Profile).
		Run()
}
