
package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"solarmesh.io/mockservices/src/service"
)


// server is used to implement protobuf.GreeterServer.
type MockTcpServer struct {
	mockService *service.MockService
}

func newMockTcpServerService(mockService *service.MockService) *MockTcpServer {
	return &MockTcpServer{mockService: mockService}
}

func init() {
	app.Register(newMockTcpServerService)
}

// Listen implementation
func (s *MockTcpServer) Listen(port string) {
	if port == "" {
		port = "8585"
	}
	tcpServerPort := fmt.Sprintf(":%v", port)
	l, err := net.Listen("tcp", tcpServerPort)
	if err != nil {
		log.Error(err)
		return
	}
	go func() {
		log.Infof("Listening on TCP port %v", port)
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			if strings.TrimSpace(string(netData)) == "STOP" {
				fmt.Println("Exiting TCP server!")
				return
			}
			// read json

			fmt.Print("-> ", string(netData))
			t := time.Now()
			myTime := t.Format(time.RFC3339) + "\n"

			// response json
			c.Write([]byte(myTime))
		}
		l.Close()
	}()

	//
	//httpHeader := make(http.Header)
	//header, ok := metadata.FromIncomingContext(ctx)
	//if ok {
	//	for k, v := range header {
	//		log.Infof("> %v: %v", k, v)
	//		httpHeader.Set(k, v[0])
	//	}
	//}
	//var resp *model.GetResponse
	//response = new(protobuf.MockResponse)
	//resp, err = s.mockService.SendRequest("TCP", nil, httpHeader)
	//if err == nil {
	//	b, e := json.Marshal(resp)
	//	if e == nil {
	//		e = json.Unmarshal(b, response)
	//	}
	//}
	//
	//// Anything linked to this variable will transmit response headers.
	//if err := ggrpc.SendHeader(ctx, header); err != nil {
	//	return nil, status.Errorf(codes.Internal, "unable to send header")
	//}

	return
}
