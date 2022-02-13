# go-micro

go-micro microservice framework template

This is a small template for creating a new microservice in go. 

Its not intended to be a fully featured microservice framework, just a small starting point with the things i normally need to build up a small, fast go microservice.

The project structure depends on this: https://github.com/golang-standards/project-layout

Features:

- usage of Opentracing/jaeger
- gelf logging
- authorization with jwt
- cached healthcheck, livez and readyz endpoints
- https/ssl and http for payload and metrics/healthcheck
- metrics with Prometheus: https://prometheus.io/docs/guides/go-application/
- Docker build with builder and target image
- chi as the router framework
- go 1.17
- automatic config substitution 

## Configuration

In this template the configuration will be automatically loaded. You have the following options to set the service configuration file.

- default: the service will try to load the configuration from the `<userhome>/<servicename>/service.yaml`
- via Commandline: `-c <configfile>` will load the configuration from this file

IN the configuration file you can use `${}` macros for adding environment variables for the configuration itself. This will not work on the `secret.yaml`. The `secret.yaml` (if given in the configuration) will load a partial configuration from another file. (Mainly for separating credentials from the other configuration) Be aware, you manually have to merge both configuration in the `config.mergeSecret()` function.



## Prometheus integration

You can switch on the prometheus integration simply by adding 

```yaml
metrics:
  enable: true
```

to the service config.

### How to add a new counter?

Simply on the class, where you want to add a new counter (or something else) make a new variable with:

```go
var (
  postConfigCounter = promauto.NewCounter(prometheus.CounterOpts{
	 Name: "gomicro_post_config_total",
     Help: "The total number of post config requests",
  })
)
```

In the code where to count the events simply do an



```go
postConfigCounter.Inc()
```

 Thats all. More examples here: https://prometheus.io/docs/guides/go-application/