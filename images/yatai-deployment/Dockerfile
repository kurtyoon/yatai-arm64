# yatai-deployment 1.1.21 for arm64: the upstream Dockerfile at the same tag,
# with GOARCH taken from the target platform instead of a fixed amd64.

FROM golang:1.27.1 AS build
ARG REPO=https://github.com/bentoml/yatai-deployment.git
# v1.1.21
ARG COMMIT=755faf260b8c010ff4d0b049405ea283e55dc4c0
ARG VERSION=1.1.21
ARG TARGETARCH
RUN git clone --quiet --filter=blob:none "$REPO" /src && git -C /src checkout --quiet "$COMMIT"
WORKDIR /src
# The same -X flags as the upstream Makefile (VERSION_BUILDFLAGS).
RUN PKG=github.com/bentoml/yatai-deployment \
    && CGO_ENABLED=0 GOOS=linux GOARCH="$TARGETARCH" GOTOOLCHAIN=local go build -a \
      -ldflags "-X '$PKG/version.GitCommit=$(echo "$COMMIT" | cut -c1-7)' -X '$PKG/version.Version=$VERSION' -X '$PKG/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
      -o /out/manager main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/manager .
USER 65532:65532
ENTRYPOINT ["/manager"]
