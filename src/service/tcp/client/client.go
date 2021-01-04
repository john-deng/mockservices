package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/silenceper/pool"
	"golang.org/x/net/context"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
	"solarmesh.io/mockservices/src/model"
)


type MockTcpClient struct {
	pool pool.Pool
}

func newMockTcpClientService() *MockTcpClient {
	return &MockTcpClient{
	}
}

func init() {
	app.Register(newMockTcpClientService)
}

func (s *MockTcpClient) connect(address string)  (err error) {
	//create a pool
	if s.pool == nil || s.pool != nil && s.pool.Len() == 0 {
		poolConfig := &pool.Config{
			InitialCap: 2,
			MaxIdle:    4,
			MaxCap:     5,
			Factory:    func() (interface{}, error) { return net.Dial("tcp", address) },
			Close:      func(v interface{}) error { return v.(net.Conn).Close() },
			// When connection reached maximum, it will close after timeout to avoid EOF issue
			IdleTimeout: 15 * time.Second,
		}
		s.pool, err = pool.NewChannelPool(poolConfig)
		if err != nil {
			log.Errorf("Error: %v", err)
		}
	}
	return
}

func (s *MockTcpClient) Send(ctx context.Context, address string, header http.Header) (response *model.TcpResponse, err error) {
	response = new(model.TcpResponse)
	var conn net.Conn

	err = s.connect(address)
	if err == nil {
		var v interface{}
		v, err = s.pool.Get()
		defer s.pool.Put(v)
		if err == nil {
			conn = v.(net.Conn)
			var b []byte
			var num int
			var resp []byte
			b, err = json.Marshal(header)
			num, err = fmt.Fprintf(conn, string(b)+"\n")
			if err == nil {
				log.Debugf("%v read", num)
				resp, err = bufio.NewReader(conn).ReadBytes('\n')
				if err == nil {
					err = json.Unmarshal(resp, response)
				}
			}
		}
	}
	return
}