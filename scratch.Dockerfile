##### BUILDER #####

# there is no need to build everything again. we already have the alpine based image and can use its atifacts 

##### TARGET #####
ARG RELEASE

FROM gomicro:alpine-${RELEASE} AS copy-src

FROM scratch

ARG RELEASE
ENV IMG_VERSION="${RELEASE}"

# hadolint ignore=DL3022
COPY --from=copy-src /usr/local/bin/gomicro /

# hadolint ignore=DL3022
COPY --from=copy-src /config/service.yaml /config/

ENTRYPOINT ["/gomicro"]
CMD ["--config","/config/service.yaml"]

EXPOSE 8080 8443

LABEL org.opencontainers.image.title="GoMicro Service" \
      org.opencontainers.image.description="GoMicro template service called GoMicro." \
      org.opencontainers.image.version="${IMG_VERSION}" \
      org.opencontainers.image.source="https://bitbucket.easy.de/scm/dm/service-gomicro-go.git" \
      org.opencontainers.image.vendor="EASY SOFTWARE AG (www.easy-software.com)" \
      org.opencontainers.image.authors="EASY SOFTWARE AG" \
      maintainer="EASY SOFTWARE AG" \
      NAME="gomicro-service"

