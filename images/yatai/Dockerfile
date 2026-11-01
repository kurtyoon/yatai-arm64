# Yatai 1.1.13 (API server and dashboard) for arm64.
# The upstream image (Dockerfile at the same tag) is alpine plus four things:
# statics/, api-server/db/ and the api-server binary from the source tree, and
# dashboard/build from the UI build. The binary is rebuilt here for the target
# CPU; the dashboard is static files, copied from the upstream amd64 image.

ARG UPSTREAM_IMAGE=quay.io/bentoml/yatai:1.1.13

FROM --platform=linux/amd64 ${UPSTREAM_IMAGE} AS upstream

FROM golang:1.27.1 AS build
ARG REPO=https://github.com/bentoml/Yatai.git
# v1.1.13
ARG COMMIT=75d516dc5df0c5bdbd13f9e1b5389a839a6e001f
ARG VERSION=1.1.13
ARG TARGETARCH
RUN git clone --quiet --filter=blob:none "$REPO" /src && git -C /src checkout --quiet "$COMMIT"
WORKDIR /src
# The same -X flags as the upstream Makefile (VERSION_BUILDFLAGS).
RUN PKG=github.com/bentoml/yatai \
    && CGO_ENABLED=0 GOOS=linux GOARCH="$TARGETARCH" GOTOOLCHAIN=local go build \
      -ldflags "-X '$PKG/api-server/version.GitCommit=$(echo "$COMMIT" | cut -c1-7)' -X '$PKG/api-server/version.Version=$VERSION' -X '$PKG/api-server/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
      -o /out/api-server ./api-server/main.go

# Upstream used alpine 3.16.2; the chart starts the server through sh.
FROM alpine:3.22
RUN mkdir -p /app/dashboard /app/scripts
COPY --from=build /src/statics /app/statics
COPY --from=upstream /app/dashboard/build /app/dashboard/build
COPY --from=build /src/api-server/db /app/db
COPY --from=build /out/api-server /app/api-server
WORKDIR /app
