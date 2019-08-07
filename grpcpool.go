package grpcpool

import (
	"google.golang.org/grpc"
	"sync"
)

func NewConnectionPool(activeCount int, dialFunc func() (*grpc.ClientConn, error)) (*ConnectionPool, error) {
	pool := &ConnectionPool{mu: sync.Mutex{}}
	for i := 0; i < activeCount; i ++ {
		client, err := dialFunc()
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.put(client)
	}
	return pool, nil
}
