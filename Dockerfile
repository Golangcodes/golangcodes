# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24
ARG ALPINE_VERSION=3.20
ARG AIR_VERSION=v1.52.3

# ---------------------------------------------------------------- dependencies
# Shared by the dev and build stages so the module download is cached once and
# only re-run when go.mod or go.sum actually change.
FROM golang:${GO_VERSION}-alpine AS deps

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

# ------------------------------------------------------------------------ dev
# Hot-reloading development server. Compose bind-mounts the repo over /src, so
# this stage holds the toolchain and nothing else.
FROM deps AS dev

ARG AIR_VERSION
RUN go install github.com/air-verse/air@${AIR_VERSION}

ENV PORT=8080
EXPOSE 8080

# Build into tmp-docker/ rather than the tmp/ that `make dev` uses on the host.
# Compose mounts a volume there, and Docker creates a bind mount's mountpoint as
# root - pointed at tmp/, that would leave a root-owned directory in the working
# tree and break the host dev server afterwards.
CMD ["air", "-c", ".air.toml", \
     "-tmp_dir", "tmp-docker", \
     "-build.bin", "./tmp-docker/main", \
     "-build.cmd", "go build -o ./tmp-docker/main ."]

# ---------------------------------------------------------------------- build
FROM deps AS build

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

# -------------------------------------------------------------------- runtime
# Code submitted from the editor is compiled by the Go playground rather than
# locally, so the image needs no Go toolchain - only the CA certificates for the
# outbound HTTPS call to play.golang.org.
FROM alpine:${ALPINE_VERSION} AS runtime

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 10001 golangcodes

WORKDIR /app

COPY --from=build /out/server ./server
# Read from disk at runtime, so they have to come along: templates are parsed
# at startup and static/ is served straight off the filesystem.
COPY templates ./templates
COPY static ./static

USER golangcodes

# Railway and most platforms inject PORT; this is the fallback.
ENV PORT=8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${PORT}/healthz" >/dev/null || exit 1

CMD ["./server"]
