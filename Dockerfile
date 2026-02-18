# ── Build stage ──────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=1 go build -ldflags="-s -w -X main.version=${VERSION}" -o /bin/agenteats-api ./cmd/api
RUN CGO_ENABLED=1 go build -ldflags="-s -w -X main.version=${VERSION}" -o /bin/agenteats-mcp ./cmd/mcp
RUN CGO_ENABLED=1 go build -ldflags="-s -w -X main.version=${VERSION}" -o /bin/agenteats-seed ./cmd/seed

# ── Runtime stage ───────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/agenteats-api /usr/local/bin/
COPY --from=builder /bin/agenteats-mcp /usr/local/bin/
COPY --from=builder /bin/agenteats-seed /usr/local/bin/

ENV HOST=0.0.0.0
ENV PORT=8000

EXPOSE 8000

ENTRYPOINT ["agenteats-api"]
