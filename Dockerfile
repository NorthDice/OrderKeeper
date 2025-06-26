# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS build

RUN apk add --no-cache git ca-certificates tzdata
RUN adduser -D -s /bin/sh -u 1001 appuser

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
COPY cmd/ ./cmd/
COPY configs/ ./configs/
COPY internal/ ./internal/
COPY migrations/ ./migrations/
COPY server/ ./server/
COPY .env .env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /order-keeper ./cmd

FROM scratch

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /order-keeper /order-keeper
COPY --from=build /app/migrations/ /migrations/
COPY --from=build /app/configs/ /app/configs/
COPY --from=build /app/.env .env

USER appuser

EXPOSE 8080

ENTRYPOINT ["/order-keeper"]