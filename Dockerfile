FROM --platform=$BUILDPLATFORM golang:1.22 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR "/app"
COPY index.html index.html
COPY server.go server.go
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -o server server.go

FROM jauderho/yt-dlp:latest
WORKDIR "/app"
COPY --from=builder /app/server .

ENTRYPOINT []
CMD "/app/server"
