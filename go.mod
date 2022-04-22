module github.com/uesleicarvalhoo/go-collector-service

go 1.18

require github.com/streadway/amqp v1.0.0

require (
	github.com/aws/aws-sdk-go v1.43.41
	github.com/netflix/go-env v0.0.0-20210215222557-e437a7e7f9fb
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/exporters/jaeger v1.6.3
	go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace v1.6.3
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)
