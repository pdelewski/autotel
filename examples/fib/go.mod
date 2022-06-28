module github.com/pdelewski/autotel/examples/fib

go 1.18

replace github.com/pdelewski/autotel => ../..

require github.com/pdelewski/autotel v0.0.0-20220627214309-9a75c355a5bd

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
)
