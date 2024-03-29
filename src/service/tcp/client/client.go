package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/john-deng/mockservices/src/model"
	"github.com/silenceper/pool"
	"golang.org/x/net/context"
)

type ConnPool map[string]pool.Pool

type MockTcpClient struct {
	connPool ConnPool
	mutex    sync.Mutex
}

func newMockTcpClientService() *MockTcpClient {
	return &MockTcpClient{
		connPool: make(ConnPool),
	}
}

func init() {
	app.Register(newMockTcpClientService)
}

func (s *MockTcpClient) connect(address string) (p pool.Pool, err error) {
	//create a pool
	p = s.connPool[address]
	if p == nil {
		poolConfig := &pool.Config{
			InitialCap: 2,
			MaxIdle:    4,
			MaxCap:     5,
			Factory:    func() (interface{}, error) { return net.Dial("tcp", address) },
			Close:      func(v interface{}) error { return v.(net.Conn).Close() },
			// When connection reached maximum, it will close after timeout to avoid EOF issue
			IdleTimeout: 15 * time.Second,
		}
		p, err = pool.NewChannelPool(poolConfig)
		if err != nil {
			log.Errorf("Error: %v", err)
			return
		}
		s.connPool[address] = p
	}

	if p != nil {
		log.Infof("TCP client connected to %v", address)
	}
	return
}

func (s *MockTcpClient) Send(ctx context.Context, address string, header http.Header) (response *model.TcpResponse, err error) {
	s.mutex.Lock()
	response = new(model.TcpResponse)
	var conn net.Conn
	var connPool pool.Pool
	connPool, err = s.connect(address)
	if err == nil {
		var v interface{}
		v, err = connPool.Get()
		defer connPool.Put(v)
		if err == nil {
			conn = v.(net.Conn)
			var b []byte
			var num int
			var resp []byte
			b, err = json.Marshal(header)
			num, err = fmt.Fprintf(conn, string(b)+"\n")
			if err == nil {
				log.Debugf("%v bytes read from TCP", num)
				resp, err = bufio.NewReader(conn).ReadBytes('\n')
				if err == nil {
					err = json.Unmarshal(resp, response)
					if err == nil {
						log.Infof("Response Header from TCP")
						for k, v := range response.Header {
							log.Infof("< %v: %v", k, v)
						}
					}
				}
			}
		}
	}
	s.mutex.Unlock()
	return
}
