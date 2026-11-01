# bentoml-proxy 0.0.1 for arm64: container_images/bentoml-proxy/Dockerfile
# from yatai-deployment v1.1.21, with the aarch_64 envoy release binary.
# yatai-deployment runs it with: envoy --config-path /etc/envoy/envoy.yaml

FROM ubuntu:22.04 AS fetch
ARG ENVOY_VERSION=1.24.1
ARG ENVOY_SHA256=e1b8a8105c72ce6d72d19c1d93ee46856815e23a540d40451426494f980f3efc
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*
RUN curl -fsSL "https://github.com/envoyproxy/envoy/releases/download/v${ENVOY_VERSION}/envoy-${ENVOY_VERSION}-linux-aarch_64" -o /usr/local/bin/envoy \
    && echo "${ENVOY_SHA256}  /usr/local/bin/envoy" | sha256sum -c - \
    && chmod +x /usr/local/bin/envoy

FROM ubuntu:22.04
COPY --from=fetch /usr/local/bin/envoy /usr/local/bin/envoy
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

# The same non-root user ID as upstream (1000).
RUN groupadd --gid 1000 envoy && useradd --uid 1000 --gid 1000 -m envoy
USER 1000:1000
