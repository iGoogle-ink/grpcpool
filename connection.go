package grpcpool

import "google.golang.org/grpc"

type Connection interface {
	Close() error
	Get() *grpc.ClientConn
}

type GrpcConnection struct {
	pool     Pool
	GrpcConn *grpc.ClientConn
}

func (this *GrpcConnection) Close() error {
	return this.pool.put(this.GrpcConn)
}

func (this *GrpcConnection) Get() *grpc.ClientConn {
	return this.GrpcConn
}
