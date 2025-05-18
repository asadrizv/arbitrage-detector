# syntax=docker/dockerfile:1
FROM golang:1.21-bullseye AS build
WORKDIR /app
COPY . .
RUN go build -o arbitrage ./cmd/server

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=build /app/arbitrage /app/arbitrage
CMD ["/app/arbitrage"]
