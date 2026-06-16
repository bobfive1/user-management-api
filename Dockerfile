# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/user-management-api ./cmd/api

FROM alpine:3.22 AS runtime

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/user-management-api /app/user-management-api

ENV GIN_MODE=release
EXPOSE 8080

USER app

CMD ["sh", "-c", "APISERVER_ADDRESS=\"0.0.0.0:${PORT:-8080}\" exec /app/user-management-api"]
