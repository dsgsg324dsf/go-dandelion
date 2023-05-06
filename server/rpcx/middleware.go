package rpcx

import (
	"context"
	"github.com/gly-hub/go-dandelion/logger"
	"github.com/gly-hub/go-dandelion/telemetry"
	"github.com/smallnest/rpcx/protocol"
)

type ClientLoggerPlugin struct {
}

func (p *ClientLoggerPlugin) PostCall(ctx context.Context, servicePath, serviceMethod string, args interface{}) error {
	requestId := logger.GetRequestId()
	ctx = context.WithValue(ctx, "request_id", requestId)
	return nil
}

type ServerLoggerPlugin struct {
}

func (p *ServerLoggerPlugin) PreHandleRequest(ctx context.Context, r *protocol.Message) error {
	logger.SetRequestId(r.Metadata["request_id"])
	traceId := r.Metadata["span_trace_id"]
	if traceId != "" {
		span, spanTraceId, err := telemetry.StartSpan(r.ServiceMethod, traceId, true)
		if err == nil {
			telemetry.SpanSetTag(span, "request_id", r.Metadata["request_id"])
			telemetry.SetSpanTraceId(spanTraceId)
			defer func() {
				telemetry.FinishSpan(span)
			}()
		}
	}

	logger.Info("client: %s, server: %v, func: %s, params: %s", r.Metadata["client_name"], r.ServicePath, r.ServiceMethod, r.Metadata["content"])
	return nil
}

func (p *ServerLoggerPlugin) PostWriteResponse(ctx context.Context, req *protocol.Message, res *protocol.Message, err error) error {
	logger.DeleteRequestId()
	return nil
}
