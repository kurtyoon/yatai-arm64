# yatai-image-builder 1.2.56 for arm64: the upstream Dockerfile at the same
# tag, with GOARCH taken from the target platform instead of a fixed amd64.
# Pin this line: an unpinned install gives 3.x, a different project line.

FROM golang:1.27.1 AS build
ARG REPO=https://github.com/bentoml/yatai-image-builder.git
# v1.2.56
ARG COMMIT=e33c6099de32ec5397cf7aff29c74c2368693c70
ARG VERSION=1.2.56
ARG TARGETARCH
RUN git clone --quiet --filter=blob:none "$REPO" /src && git -C /src checkout --quiet "$COMMIT"
WORKDIR /src
# The same -X flags as the upstream Makefile (VERSION_BUILDFLAGS).
RUN PKG=github.com/bentoml/yatai-image-builder \
    && CGO_ENABLED=0 GOOS=linux GOARCH="$TARGETARCH" GOTOOLCHAIN=local go build -a \
      -ldflags "-X '$PKG/version.GitCommit=$(echo "$COMMIT" | cut -c1-7)' -X '$PKG/version.Version=$VERSION' -X '$PKG/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
      -o /out/manager main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/manager .
USER 65532:65532
ENTRYPOINT ["/manager"]
