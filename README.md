Here’s a tight, senior-friendly `README.md` you can drop in and tweak.

---

# Your Framework (batteries-included Go toolkit)

A pragmatic framework to speed up backend development with sane defaults and **batteries included**:

* HTTP endpoint service
* gRPC endpoint service
* Message broker (RMQ) producer/consumer
* Monitoring (Prometheus metrics, health probes)
* Pyroscope profiling (continuous)

## Table of Contents

* [Why](#why)
* [Features](#features)
* [Quick Start](#quick-start)
* [Install](#install)
* [Project Layout](#project-layout)
* [Configuration](#configuration)
* [HTTP](#http)
* [gRPC](#grpc)
* [RabbitMQ](#rabbitmq)
* [Monitoring & Observability](#monitoring--observability)
* [Pyroscope](#pyroscope)
* [Local Dev Stack](#local-dev-stack)
* [Make/Task Targets](#maketask-targets)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)

## Why

Stop re-wiring the same plumbing every project. This repo gives you **production-grade defaults** out of the box: servers, middlewares, metrics, tracing hooks, graceful shutdown, and message queues — so you ship business logic, not boilerplate.

## Features

* **HTTP**: idiomatic router, request validation, structured logging, graceful shutdown.
* **gRPC**: reflection, health service, interceptors (logging/metrics/recovery).
* **RMQ**: publisher/consumer with reconnection, prefetch/QoS, DLQ pattern.
* **Monitoring**: `/metrics`, `/healthz`, `/readyz`, Prometheus counters/histograms.
* **Profiling**: Pyroscope agent wired by env vars; opt-in in production.
* **12-factor config**: env-first with YAML override.
* **Zero surprises**: context-aware, cancellation-safe, backoff & retries where it matters.

## Quick Start

```bash
# 1) clone
git clone https://github.com/your-org/yourfwk.git
cd yourfwk

# 2) run minimal service (http+grpc+metrics)
go run ./cmd/app
```

Hit:

* HTTP: [http://localhost:8080/hello](http://localhost:8080/hello)
* Metrics: [http://localhost:8080/metrics](http://localhost:8080/metrics)
* Health:  [http://localhost:8080/healthz](http://localhost:8080/healthz)
* gRPC:    localhost:9090 (with reflection)

## Install

Use as a library in another service:

```bash
go get github.com/your-org/yourfwk@latest
```

Minimal `main.go`:

```go
package main

import (
	"context"
	"log"

	fwk "github.com/your-org/yourfwk"
)

func main() {
	app := fwk.NewApp(
		fwk.WithHTTP(),
		fwk.WithGRPC(),
		fwk.WithRabbitMQ(),
		fwk.WithMonitoring(),
		fwk.WithPyroscope(),
	)
	// Register routes, gRPC servers, consumers here
	app.HTTP.GET("/hello", func(c fwk.Context) error { return c.JSON(200, fwk.M{"ok": true}) })
	app.Run(context.Background())
	log.Println("bye")
}
```

## Project Layout

```
.
├─ cmd/
│  └─ app/                 # main entry
├─ internal/
│  ├─ http/                # http handlers, dto, middleware
│  ├─ grpc/                # grpc services & interceptors
│  ├─ mq/                  # RMQ publishers/consumers
│  ├─ metrics/             # custom metrics
│  └─ config/              # config loader & defaults
├─ pkg/                    # reusable helpers
├─ deploy/
│  └─ docker-compose/      # rabbitmq, prometheus, grafana, pyroscope
└─ README.md
```

## Configuration

Environment-first; YAML optional. Key envs:

```env
# servers
HTTP_ADDR=:8080
GRPC_ADDR=:9090
SHUTDOWN_TIMEOUT=15s

# logging & metrics
LOG_LEVEL=info
METRICS_ENABLED=true

# rabbitmq
RMQ_URL=amqp://guest:guest@localhost:5672/
RMQ_PREFETCH=32
RMQ_PUBLISHER_CONFIRM=true
RMQ_EXCHANGE=app.topic
RMQ_DLX=app.dlx

# pyroscope
PYROSCOPE_ENABLED=true
PYROSCOPE_SERVER=http://localhost:4040
PYROSCOPE_APP_NAME=yourfwk-app
PYROSCOPE_PROFILE_TYPES=cpu,alloc,inuse
```

Optional `config.yaml`:

```yaml
http:
  addr: ":8080"
grpc:
  addr: ":9090"
rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
```

## HTTP

Register endpoints with common middleware (request ID, recovery, metrics):

```go
app.HTTP.Use(fwk.Middleware.RequestID(), fwk.Middleware.Recover(), fwk.Middleware.Metrics())

app.HTTP.GET("/hello", func(c fwk.Context) error {
	type Resp struct{ Message string `json:"message"` }
	return c.JSON(200, Resp{Message: "hello"})
})
```

Health/Ready endpoints are auto-mounted at `/healthz` and `/readyz`.

## gRPC

Unary/stream interceptors (logging, metrics, recovery) + reflection + health:

```go
s := app.GRPC.Server()
pb.RegisterGreeterServer(s, &Greeter{}) // your implementation
```

Client example lives under `examples/grpc-client/`.

## RabbitMQ

Producer:

```go
pub := app.RMQ.Publisher(fwk.RMQPublisherConfig{
	Exchange: "app.topic",
	Kind:     "topic",
	Confirm:  true,
})
_ = pub.Publish(ctx, "user.created", fwk.ContentJSON(map[string]any{"id": 42}))
```

Consumer:

```go
app.RMQ.Consumer(fwk.RMQConsumerConfig{
	Queue:       "user.created.q",
	Bind:        []fwk.Bind{{Exchange: "app.topic", RoutingKey: "user.created"}},
	Prefetch:    32,
	AutoAck:     false,
	DeadLetter:  fwk.DeadLetter{Exchange: "app.dlx"},
	Concurrency: 4,
}, func(ctx context.Context, d fwk.Delivery) error {
	// handle message
	return d.Ack()
})
```

**Notes**

* Publisher confirms enabled by default in prod.
* Consumers auto-reconnect with exponential backoff.
* DLX recommended; see `deploy/docker-compose/rabbitmq.conf`.

## Monitoring & Observability

* **Prometheus metrics** at `/metrics` (HTTP).
* Default metrics: HTTP/gRPC duration histograms, RMQ delivery counters, Go runtime.
* Add your own:

```go
var productCreateLatency = fwk.Metrics.NewHistogramVec("product_create_seconds", "Latency", []string{"result"})
```

* **Health/Readiness** endpoints:

  * `/healthz`: liveness
  * `/readyz`: checks RMQ channel, port binds, etc.

## Pyroscope

Enabled via env. The agent starts with the app and pushes profiles.

```go
// on app init (if PYROSCOPE_ENABLED=true)
fwk.Pyroscope.Start(fwk.PyroscopeConfig{
	Server:   os.Getenv("PYROSCOPE_SERVER"),
	AppName:  os.Getenv("PYROSCOPE_APP_NAME"),
	Profiles: strings.Split(os.Getenv("PYROSCOPE_PROFILE_TYPES"), ","),
})
```

Open Pyroscope UI → explore `cpu`, `alloc_bytes`, `inuse_objects`, etc.

## Local Dev Stack

`deploy/docker-compose/docker-compose.yml` spins up RabbitMQ, Prometheus, Grafana, Pyroscope.

```yaml
version: "3.9"
services:
  rabbitmq:
    image: rabbitmq:3.13-management
    ports: ["5672:5672", "15672:15672"]
  prometheus:
    image: prom/prometheus
    volumes: ["./prometheus.yml:/etc/prometheus/prometheus.yml"]
    ports: ["9091:9090"]
  grafana:
    image: grafana/grafana
    ports: ["3000:3000"]
  pyroscope:
    image: grafana/pyroscope:latest
    ports: ["4040:4040"]
```

Run:

```bash
docker compose -f deploy/docker-compose/docker-compose.yml up -d
```

Prometheus scrapes `http://host.docker.internal:8080/metrics` by default (adjust for Linux).

## Make/Task Targets

Use either `make` or `task`:

```makefile
run:        ## Run app
	go run ./cmd/app
lint:       ## Lint
	golangci-lint run
test:       ## Unit tests
	go test ./...
```

## Roadmap

* [ ] Kafka provider (in addition to RMQ)
* [ ] OpenTelemetry traces/logs wiring
* [ ] Graceful rolling restart hooks
* [ ] Config hot-reload
* [ ] CLI scaffolding (`yourfwk new service`)

## Contributing

PRs welcome. Keep code small, focused, and covered. Follow existing patterns:

* Context everywhere, no global singletons
* Return typed errors; wrap with `%w`
* Add metrics around external I/O

## License

MIT (or your choice). See `LICENSE`.

---

**Tip:** keep example services under `examples/` so new users can copy-paste a working HTTP+gRPC+RMQ starter in minutes.
