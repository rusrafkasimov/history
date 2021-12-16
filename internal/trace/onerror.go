package trace

import (
	"context"
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/rusrafkasimov/history/internal/config"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

func MakeSpan(ctx context.Context, tracer opentracing.Tracer, opnName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}

	span := tracer.StartSpan(opnName, opts...)

	return span
}

func OnError(logger promtail.Client, span opentracing.Span, err error) {
	if span != nil {
		span.SetTag(string(ext.Error), true)
		span.LogKV(otlog.Error(err))
		span.LogFields(otlog.String("error", err.Error()))
	}
	logger.Errorf("%v", err.Error())
}

func InitJaegerTracing(ctx context.Context, contextKeyName interface{}, configuration *config.Configuration) (closer io.Closer, err error) {
	host, err := configuration.Get("JAEGER_AGENT_HOST")
	if err != nil {
		return nil, err
	}

	port, err := configuration.Get("JAEGER_AGENT_PORT")
	if err != nil {
		return nil, err
	}

	cfg := jaegercfg.Configuration{
		ServiceName: ctx.Value(contextKeyName).(string),
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%s", host, port),
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, _ := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
