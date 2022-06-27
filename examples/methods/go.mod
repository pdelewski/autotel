module github.com/pdelewski/autotel/examples/methods

go 1.18

replace github.com/pdelewski/autotel => ../..

replace github.com/pdelewski/autotel/rtlib => ../../rtlib

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel v1.6.3 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.6.3 // indirect
	go.opentelemetry.io/otel/sdk v1.6.3 // indirect
	go.opentelemetry.io/otel/trace v1.6.3 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
	github.com/pdelewski/autotel/rtlib v0.0.0-00010101000000-000000000000 // indirect
)
