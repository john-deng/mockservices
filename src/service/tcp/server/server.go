package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"solarmesh.io/mockservices/src/model"
	"solarmesh.io/mockservices/src/service"
)

// server is used to implement protobuf.GreeterServer.
type MockServer struct {
	mockService *service.MockService
}

func newMockTcpServerService(mockService *service.MockService) *MockServer {
	return &MockServer{mockService: mockService}
}

func init() {
	app.Register(newMockTcpServerService)
}

// Listen implementation
func (s *MockServer) Listen(port string) {
	go func() {
		for {
			var err error
			if port == "" {
				port = "8585"
			}
			address := fmt.Sprintf(":%v", port)
			l, err := net.Listen("tcp", address)
			if err != nil {
				log.Error(err)
				return
			}
			log.Infof("Listening on TCP port %v", port)
			var conn net.Conn
			conn, err = l.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			var request *model.TcpRequest
			var response *model.TcpResponse
			for {
				netData, err := bufio.NewReader(conn).ReadBytes('\n')
				if err != nil {
					log.Info(err)
					break
				}
				if strings.TrimSpace(string(netData)) == "STOP" {
					log.Info("Exiting TCP server!")
					break
				}

				var rb []byte
				var num int
				request = new(model.TcpRequest)
				request.Header = make(http.Header)
				response = new(model.TcpResponse)
				// parse data
				//log.Debugf("request: %v", string(netData))
				err = json.Unmarshal(netData, &request.Header)
				if err == nil {
					response.Response, err = s.mockService.SendRequest("TCP", nil, request.Header)
					response.Header = request.Header
					rb, err = json.Marshal(response)
					rb = append(rb, []byte("\n")...)
					if err == nil {
						num, err = conn.Write(rb)
						if err == nil {
							log.Infof("%v bytes response by TCP", num)
						}
					}
				}
			}
			err = l.Close()
		}
	}()
}
