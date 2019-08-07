package grpcpool

import (
	"google.golang.org/grpc"
	"sync"
	"container/list"
	"errors"
)

type Pool interface {
	Get() (Connection, error)
	Close()
	put(client *grpc.ClientConn) error
}

type ConnectionPool struct {
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	idle   list.List
}

func (this *ConnectionPool) Get() (Connection, error) {
	conn, err := this.get()
	if err != nil {
		return nil, err
	}
	return &GrpcConnection{GrpcConn: conn, pool: this}, nil
}

func (this *ConnectionPool) Close() {
	this.mu.Lock()
	idle := this.idle
	this.idle.Init()
	this.closed = true
	if this.cond != nil {
		this.cond.Broadcast()
	}
	this.mu.Unlock()
	for e := idle.Front(); e != nil; e = e.Next() {
		e.Value.(*grpc.ClientConn).Close()
	}
}

func (this *ConnectionPool) put(client *grpc.ClientConn) error {
	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return client.Close()
	}

	this.idle.PushFront(client)
	this.mu.Unlock()
	if this.cond != nil {
		this.cond.Signal()
	}
	return nil
}

func (this *ConnectionPool) get() (*grpc.ClientConn, error) {
	this.mu.Lock()
	for {
		element := this.idle.Front()

		if element != nil {
			this.idle.Remove(element)
			client := element.Value.(*grpc.ClientConn)
			this.mu.Unlock()
			return client, nil
		}

		if this.closed {
			return nil, errors.New("Pool is closed")
		}

		if this.cond == nil {
			this.cond = sync.NewCond(&this.mu)
		}
		this.cond.Wait()
	}

}
