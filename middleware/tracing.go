package middleware

import (
	"fmt"
	opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go-micro.dev/v4/server"
	"golang.org/x/exp/slices"
	"golang.org/x/net/context"
)

func NewTracingHandlerWrapper(ot opentracing.Tracer) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			var err error
			if ot == nil {
				ot = opentracing.GlobalTracer()
			}
			path := req.Endpoint()
			if slices.Contains(skipPaths, path) {
				err = h(ctx, req, rsp)
				return err
			}
			name := fmt.Sprintf("%s.%s", req.Service(), path)
			ctx, span, err := opentracing4.StartSpanFromContext(ctx, ot, name)
			if err != nil {
				return err
			}
			defer span.Finish()
			if err = h(ctx, req, rsp); err != nil {
				span.LogFields(opentracinglog.String("error", err.Error()))
				span.SetTag("error", true)
			}
			return err
		}
	}
}
