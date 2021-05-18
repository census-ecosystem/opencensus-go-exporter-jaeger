module contrib.go.opencensus.io/exporter/jaeger

go 1.16

require (
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	go.opencensus.io v0.22.4
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	google.golang.org/api v0.29.0
)

replace github.com/uber/jaeger-client-go v2.25.0+incompatible => github.com/nhatthm/jaeger-client-go v2.28.1-0.20210518145049-c75020212e9e+incompatible
