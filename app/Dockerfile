FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o orglang app/main.go

FROM alpine

WORKDIR /app

COPY app/reference.yaml .

COPY --from=builder /build/orglang .

ENTRYPOINT ["/app/orglang"]
