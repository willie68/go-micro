##### BUILDER #####

FROM golang:1.13-alpine3.11 as builder

## Task: Install build deps

# hadolint ignore=DL3018
RUN set -eux; \
    apk add --no-progress --quiet --no-cache --upgrade --virtual .build-deps \
        gcc \
        git \
        musl-dev

## Task: copy source files

COPY . /src
WORKDIR /src

## Task: fetch project deps

RUN go mod download

## Task: build project

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV CGO_ENABLED="0"

RUN go build -ldflags="-s -w" -o gomicro cmd/service.go 

## Task: set permissions

RUN chmod 0755 /src/gomicro

## Task: runtime dependencies

# hadolint ignore=DL3018
RUN set -eux; \
    apk add --no-progress --quiet --no-cache --upgrade --virtual .run-deps \
        tzdata

# hadolint ignore=DL3018,SC2183,DL4006
RUN set -eu +x; \
    apk add --no-progress --quiet --no-cache --upgrade ncurses; \
    apk update --quiet; \
    printf '%30s\n' | tr ' ' -; \
    echo "RUNTIME DEPENDENCIES"; \
    PKGNAME=$(apk info --depends .run-deps \
        | sed '/^$/d;/depends/d' \
        | sort -u ); \
    printf '%s\n' "${PKGNAME}" \
        | while IFS= read -r pkg; do \
                apk info --quiet --description --no-network "${pkg}" \
                | sed -n '/description/p' \
                | sed -r "s/($(echo "${pkg}" | sed -r 's/\+/\\+/g'))-(.*)\s.*/\1=\2/"; \
                done \
        | tee -a /usr/share/rundeps; \
    printf '%30s\n' | tr ' ' - 


##### TARGET #####

FROM alpine:3.11

ARG RELEASE
ENV IMG_VERSION="${RELEASE}"

COPY --from=builder /src/gomicro /usr/local/bin/
COPY --from=builder /src/configs/service.yaml /config/
COPY --from=builder /usr/share/rundeps /usr/share/rundeps

RUN set -eux; \
    xargs -a /usr/share/rundeps apk add --no-progress --quiet --no-cache --upgrade --virtual .run-deps

ENTRYPOINT ["/usr/local/bin/gomicro"]
CMD ["--config","/config/service.yaml"]

EXPOSE 8080 8443

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=10s \
  CMD wget -q -T 5 --spider http://localhost:8080/health/health

LABEL org.opencontainers.image.title="GoMicro" \
      org.opencontainers.image.description="DM GoMicro" \
      org.opencontainers.image.version="${IMG_VERSION}" \
      org.opencontainers.image.source="https://bitbucket.easy.de/scm/dm/service-gomicro-go.git" \
      org.opencontainers.image.vendor="EASY SOFTWARE AG (www.easy-software.com)" \
      org.opencontainers.image.authors="EASY Apiomat GmbH" \
      maintainer="EASY Apiomat GmbH" \
      NAME="gomicro"

