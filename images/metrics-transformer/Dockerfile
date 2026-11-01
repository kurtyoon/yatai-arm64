FROM golang:1.27.1 AS build
ARG TARGETARCH
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH="$TARGETARCH" go build -trimpath -ldflags "-s -w" -o /out/metrics-transformer .

# yatai-deployment sets no command for this container, so the entrypoint runs.
FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/metrics-transformer /metrics-transformer
USER 65532:65532
ENTRYPOINT ["/metrics-transformer"]
