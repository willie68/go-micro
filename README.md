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
- go 1.24
- automatic config substitution 

## Why using this and not a framework?

Because you gain more flexibility. See this little repo as a starting point for writing your own microservice framework for you or your company.  

## Configuration

In this template the configuration will be automatically loaded. You have the following options to set the service configuration file.

- default: the service will try to load the configuration from the `<userhome>/<servicename>/service.yaml`
- via Commandline: `-c <configfile>` will load the configuration from this file

In the configuration file you can use `${}` macros for adding environment variables for the configuration itself. 

### Secrets

The configuration can be split into two parts: one for the normal configuration. Additionally, certain parts, such as credentials, can be stored in a second file. (Without the possibility to use ${} macros.) These can then be made available via another mechanism (e.g., via the Kubernetes secret store). The structure must be identical. This second file is then referenced via the "secretfile" entry in the config.yaml file. When the service starts, first the config.yaml will be loaded, the macros will be processed and than the secret file (if present) will be loaded into the config.

```yaml
secretfile: "./config/secret.yaml"
```

### Enviroment

Values of config files can be replaced with environment variables. To do this, ${name} in the string is replaced with the corresponding values of the current environment variables. References to undefined variables are replaced with an empty string. (see [drone/envsubst](https://github.com/drone/envsubst))

### Prometheus integration

You can switch on the prometheus integration simply by adding 

```yaml
metrics:
  enable: true
```

to the service config.

#### How to add a new counter?

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

 That's all. More examples here: https://prometheus.io/docs/guides/go-application/