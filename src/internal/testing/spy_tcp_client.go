package testing

import (
	"sync"

	"github.com/cloudfoundry/metric-store-release/src/pkg/rpc"
	"github.com/niubaoshu/gotiny"
)

type SpyTCPClient struct {
	mu             sync.Mutex
	sentPointsChan chan *rpc.Point
	sentPoints     []*rpc.Point
	writeErrorChan chan error
}

func NewSpyTCPClient() *SpyTCPClient {
	return &SpyTCPClient{
		sentPointsChan: make(chan *rpc.Point, 1000),
		writeErrorChan: make(chan error, 10),
	}
}

func (s *SpyTCPClient) Write(data []byte) (int, error) {
	select {
	case err := <-s.writeErrorChan:
		return 0, err
	default:
	}

	size := len(data)

	batch := rpc.Batch{}
	gotiny.Unmarshal(data, &batch)

	for _, point := range batch.Points {
		s.sentPointsChan <- point
	}

	return size, nil
}

func (s *SpyTCPClient) SetErr(err error) {
	s.writeErrorChan <- err
}

func (s *SpyTCPClient) GetPoints() []*rpc.Point {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		select {
		case point := <-s.sentPointsChan:
			s.sentPoints = append(s.sentPoints, point)
		default:
			return s.sentPoints
		}
	}
}
