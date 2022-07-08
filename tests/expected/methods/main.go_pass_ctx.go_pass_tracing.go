package main

import (
	"github.com/pdelewski/autotel/rtlib"
	otel "go.opentelemetry.io/otel"
	"context"
)

type element struct {
}

type driver struct {
	e element
}

type i interface {
	foo(p int, __tracing_ctx context.Context) int
}

type impl struct {
}

func (i impl) foo(p int, __tracing_ctx context.Context) int {
	__child_tracing_ctx, span := otel.Tracer("foo").Start(__tracing_ctx, "foo")
	_ = __child_tracing_ctx
	defer span.End()
	return 5
}

func foo(p int, __tracing_ctx context.Context) int {
	__child_tracing_ctx, span := otel.Tracer("foo").Start(__tracing_ctx, "foo")
	_ = __child_tracing_ctx
	defer span.End()
	return 1
}

func (d driver) process(a int, __tracing_ctx context.Context) {
	__child_tracing_ctx, span := otel.Tracer("process").Start(__tracing_ctx, "process")
	_ = __child_tracing_ctx
	defer span.End()
}

func (e element) get(a int, __tracing_ctx context.Context) {
	__child_tracing_ctx, span := otel.Tracer("get").Start(__tracing_ctx, "get")
	_ = __child_tracing_ctx
	defer span.End()
}

func main() {
	__child_tracing_ctx := context.TODO()
	_ = __child_tracing_ctx
	ts := rtlib.NewTracingState()
	defer rtlib.Shutdown(ts)
	otel.SetTracerProvider(ts.Tp)
	ctx := context.Background()
	__child_tracing_ctx, span := otel.Tracer("main").Start(ctx, "main")
	defer span.End()
	rtlib.AutotelEntryPoint__()
	d := driver{}
	d.process(10, __child_tracing_ctx)
	d.e.get(5, __child_tracing_ctx)
	var in i
	in = impl{}
	in.foo(10, __child_tracing_ctx)
}
