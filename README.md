# go-micro

go-micro microservice framework template

This is a small template for creating a new microservice in go.

Its not intended to be a fully featured microservice framework, just a small starting point with the things i normally need to build up a small, fast go microservice.

The layout follows Clean Architecture with hexagonal inbound/outbound adapters ([Do Digitals](https://dodigitals.org/blog/golang-microservices-folder-structure-do-digitals), [Medium](https://medium.com/@gitesky14/production-ready-go-folder-structure-88c1bd0f5a07)):

```
cmd/service                     service entry point
cmd/test                        small test helper
internal/
  adapter/inbound/http/         REST routes, JWT, health mount, Swagger UI
  adapter/outbound/address/     persistence: in-memory (`internal`) or MySQL
  domain/addresses              address use case and storage port
  infrastructure/               HTTP/TLS server, health, slog logging
  bootstrap/                    samber/do wiring
  config/                       YAML load, envsubst, secret merge
  shared/                       errors, HTTP helpers, TTL cache
pkg/client                      Go client for the address API
pkg/pmodel                      public models
pkg/web                         embedded web UI
api/                            generated OpenAPI/Swagger (`swag`)
bruno/                          Bruno collection (health HTTP/HTTPS + CRUD)
docs/api                        how to regenerate Swagger
configs/                        example service and secret YAML
scripts/                        build, start, test, race, lint, docker
```

Dependency injection uses [samber/do](https://github.com/samber/do) v2. The HTTP stack is [chi](https://github.com/go-chi/chi).

Features:

- structured logging with `log/slog` (stdout, optional rolling file)
- GELF to Graylog (UDP or TCP)
- VictoriaLogs JSON-line sender (async queue; a down collector does not block the process)
- OpenTelemetry (OTLP HTTP traces)
- optional JWT authentication
- cached health checks, `/livez` and `/readyz` (GET and HEAD)
- HTTPS for the API plus HTTP for probes/metrics when TLS is enabled
- Prometheus metrics: https://prometheus.io/docs/guides/go-application/
- Docker multi-stage image
- Go 1.26
- config `${}` substitution and optional secret file merge

## Why using this and not a framework?

Because you gain more flexibility. See this little repo as a starting point for writing your own microservice framework for you or your company.

## Run locally

```
go build -o gomicro-service.exe ./cmd/service
gomicro-service.exe -c ./configs/service_local.yaml
```

Or `scripts\start.cmd` after a build. Default local ports are HTTP **9480** and HTTPS **9443**.

With TLS enabled the HTTP port only serves health, metrics and profiling. Address CRUD lives on HTTPS under `/api/v1/addresses`. All address calls need the `tenant` header. JWT is commented out in the local configs (`auth.type: #jwt`).

The TLS certificate is generated at runtime unless `http.certificate` and `http.key` are set.

## Configuration

The service loads its YAML automatically:

- default: `<userhome>/<servicename>/service/service.yaml` (`${configdir}` is the per-user config folder)
- command line: `-c <configfile>`

`${name}` in the config is replaced from the process environment ([drone/envsubst](https://github.com/drone/envsubst)). Undefined variables become an empty string.

### Secrets

Credentials can live in a second file with the same structure (no `${}` macros). Point to it with `secretfile`. The main file is loaded and substituted first, then the secret file is merged on top. Typical source is a Kubernetes secret mount.

```yaml
secretfile: "./config/secret.yaml"
```

### Logging

```yaml
logging:
  level: INFO
  filename: logging.log          # optional lumberjack file
  gelf-url:                      # Graylog host; empty disables GELF
  gelf-port: 12201
  gelf-protocol: udp             # udp (default) or tcp
  victoria-logs-url:             # e.g. http://localhost:9428; empty disables it
```

VictoriaLogs is written asynchronously. If the collector is down, records stay in a bounded in-memory queue and the logger does not hang.

### Address storage

```yaml
addressstorage:
  type: "internal"    # in-memory demo store
  # type: "mysql"
  # connection:
  #   host: 127.0.0.1
  #   database: gomicro
  #   table: addresses
  #   username: ...
  #   password: ...
```

See `configs/service_mysql.yaml` for a MySQL example.

### Prometheus

```yaml
metrics:
  enable: true
```

Add a counter where you need it:

```go
var (
  postAdrCounter = promauto.NewCounter(prometheus.CounterOpts{
    Name: "gomicro_post_adr_total",
    Help: "The total number of address requests",
  })
)

postAdrCounter.Inc()
```

More examples: https://prometheus.io/docs/guides/go-application/

## API

- OpenAPI contract: `api/swagger.yaml` (also served at `/swagger/` when the service runs)
- Regenerate with `scripts\build.cmd` or the `swag init` command in [docs/api/README.md](docs/api/README.md)
- Bruno collection: open `bruno/` in Bruno, select environment `local`. Collection variables live in `bruno/collection.bru`; environments only override ports. Use `--insecure` (or disable TLS verify) for the generated certificate. Details: [bruno/README.md](bruno/README.md)
- Go client: `pkg/client`

## Tests and scripts

| Script | Purpose |
| --- | --- |
| `scripts\unittests.cmd` | `go test -coverprofile=cover.out -coverpkg=./... ./...` then `go tool cover -func` |
| `scripts\racetest.cmd` | `set CC=clang` then `go test --race ./...` (required on Windows) |
| `scripts\lint.cmd` | `revive` with `revive.toml` |
| `scripts\build.cmd` | Swagger + `go build` |
| `scripts\dockerbuild.cmd` | image `mcs/gomicro-service:V1`, ports 9080/9543 |

## Docker

```
docker build -f ./build/package/Dockerfile ./ -t mcs/gomicro-service:V1
```

The image exposes 8080 (HTTP health) and 8443 (HTTPS API). The container health check hits `http://localhost:8080/livez`.
