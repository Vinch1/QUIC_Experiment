FROM golang:1.26 AS builder

WORKDIR /src
COPY . .
RUN go build -o /out/node ./cmd/node

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /out/node /app/node
ENTRYPOINT ["/app/node"]

