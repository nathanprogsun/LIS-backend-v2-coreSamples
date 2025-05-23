package common

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const ServiceName = "coresamples-v2"

func UnaryClientInterceptor(trackingId string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		// identifier meta data
		md := metadata.New(map[string]string{
			"service-name": ServiceName,
			"x-request-id": trackingId,
			"ip":           Env.PodIP,
			"hostname":     Env.ServiceName,
		})

		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			Errorf("gRPC call failed", err)
		}
		return err
	}
}
