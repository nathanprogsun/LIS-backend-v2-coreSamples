package middleware

import (
	"bytes"
	"context"
	"coresamples/common"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/metadata"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type Fn func(c *gin.Context) []zapcore.Field

// Config is config setting for Ginzap
type Config struct {
	TimeFormat string
	UTC        bool
	SkipPaths  []string
	Context    Fn
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

var skipPaths = []string{"Health.Check"}

var timeFormat string

var LogHandler = func(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fields := []zapcore.Field{}
		start := time.Now()
		path := req.Endpoint()
		md, ok := metadata.FromIncomingContext(ctx)
		requestIDs := md.Get("x-request-id")
		serviceNames := md.Get("service-name")
		// create request ID if not exist in the header
		var requestIdField string
		var serviceNameField string
		if len(serviceNames) == 0 {
			serviceNameField = "unknown caller"
		} else {
			serviceNameField = serviceNames[0]
		}

		if len(requestIDs) > 0 {
			requestIdField = requestIDs[0]
		} else {
			requestIdField = uuid.New().String()
		}
		// add request ID to response header
		resp, ok := rsp.(server.Response)
		if ok {
			resp.WriteHeader(map[string]string{"x-request-id": requestIdField})
			resp.WriteHeader(map[string]string{"service-name": serviceNameField})
		}
		// log request ID
		fields = append(fields, zap.String("x-request-id", requestIdField))
		// log caller service name
		fields = append(fields, zap.String("caller-service-name", serviceNameField))

		// log trace and span ID
		if trace.SpanFromContext(ctx).SpanContext().IsValid() {
			fields = append(fields, zap.String("trace_id", trace.SpanFromContext(ctx).SpanContext().TraceID().String()))
			fields = append(fields, zap.String("span_id", trace.SpanFromContext(ctx).SpanContext().SpanID().String()))
		}

		// log request body
		body, _ := json.Marshal(req.Body())
		fields = append(fields, zap.String("request_body", string(body)))

		// call next handler, create chaining
		err := fn(ctx, req, rsp)

		if err != nil {
			log.Error()
			return err
		}

		// if the request endpoint is not skip-able
		if !slices.Contains(skipPaths, path) {
			end := time.Now()
			latency := end.Sub(start)
			end = end.UTC()
			resp, _ := json.Marshal(rsp)
			fields = append(fields, []zapcore.Field{
				zap.String("method", req.Method()),
				zap.String("path", path),
				zap.String("resp_body", string(resp)),
				zap.Duration("latency", latency),
			}...)
			if timeFormat != "" {
				fields = append(fields, zap.String("time", end.Format(timeFormat)))
			}
			common.InfoFields(path, fields...)
		}

		return err
	}
}

func defaultHandleRecovery(c *gin.Context, err interface{}) {
	c.AbortWithStatus(http.StatusInternalServerError)
}

/*
RecoveryWithZap returns a gin.HandlerFunc (middleware)
that recovers from any panics and logs requests using uber-go/zap.
All errors are logged using zap.Error().
stack means whether output the stack info.
The stack info is easy to find where the error occurs but the stack info is too large.
*/
func RecoveryWithZap(fn server.HandlerFunc) server.HandlerFunc {
	//err := sentry.Init(sentry.ClientOptions{
	//	// Either set your DSN here or set the SENTRY_DSN environment variable.
	//	Dsn: "https://fde61d937fcf4562bc7de1e439be1f21@sentry1.vibrant-america.com/37",
	//	// Either set environment and release here or set the SENTRY_ENVIRONMENT
	//	// and SENTRY_RELEASE environment variables.
	//	// Environment: "Production",
	//	// Release:     "my-project-name@1.0.0",
	//	// Enable printing of SDK debug messages.
	//	// Useful when getting started or trying to figure something out.
	//	Debug:       true,
	//	Environment: os.Getenv("CORESAMPLES_ENV"),
	//	// Set TracesSampleRate to 1.0 to capture 100%
	//	// of transactions for performance monitoring.
	//	// We recommend adjusting this value in production,
	//	TracesSampleRate: 1,
	//	AttachStacktrace: true,
	//})
	//
	//if err != nil {
	//	log.Fatal("sentry.Init: %s", zap.Error(err))
	//}
	//defer sentry.Flush(2 * time.Second)

	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if err := recover(); err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 2)

			log.Error("[Recovery from panic]",
				zap.Time("time", time.Now()),
				zap.Any("error", err),
				zap.String("stack", string(debug.Stack())))
		}
		err := fn(ctx, req, rsp)
		return err
	}
}

// CustomRecoveryWithZap returns a gin.HandlerFunc (middleware) with a custom recovery handler
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func CustomRecoveryWithZap(logger *zap.Logger, stack bool, recovery gin.RecoveryFunc) gin.HandlerFunc {
	err := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: "https://fde61d937fcf4562bc7de1e439be1f21@sentry1.vibrant-america.com/37",
		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		// Environment: "Production",
		// Release:     "my-project-name@1.0.0",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: false,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 0.5,
	})

	if err != nil {
		log.Fatal("sentry.Init: %s", zap.Error(err))
	}
	defer sentry.Flush(2 * time.Second)

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 2)

				// check if the error is broken pipe
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				recovery(c, err)
			}
		}()
		c.Next()
	}
}

func Ginzap(timeFormat string, utc bool) gin.HandlerFunc {
	return GinzapWithConfig(&Config{TimeFormat: timeFormat, UTC: utc})
}

func GinzapWithConfig(conf *Config) gin.HandlerFunc {
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		fields := []zapcore.Field{}
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		requestID := c.Request.Header.Get("X-Request-Id")
		var requestIdField string
		if requestID != "" {
			requestIdField = requestID
		} else {
			requestIdField = uuid.New().String()
		}

		c.Writer.Header().Set("X-Request-Id", requestIdField)
		fields = append(fields, zap.String("request_id", requestIdField))
		// log trace and span ID
		if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
			fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
			fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
		}

		// log request body
		var body []byte
		var buf bytes.Buffer

		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ = io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		fields = append(fields, zap.String("body", string(body)))

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}
			accountID := c.GetString("account_id")
			accountType := c.GetString("account_type")
			fields = append(fields, []zapcore.Field{
				zap.Int("status", c.Writer.Status()),
				zap.String("resp_body", w.body.String()),
				zap.String("account_id", accountID),
				zap.String("account_type", accountType),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
			}...)
			if conf.TimeFormat != "" {
				fields = append(fields, zap.String("time", end.Format(conf.TimeFormat)))
			}

			if conf.Context != nil {
				fields = append(fields, conf.Context(c)...)
			}
			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					common.ErrorFields(e, fields...)
				}
			} else {
				common.InfoFields(path, fields...)
			}
		}
	}
}
