# Builder
FROM golang:1.16 as builder
LABEL stage=builder

WORKDIR /builder

COPY . .

RUN CGO_ENABLED=0 go build -o ./bin/chart-viewer ./cmd/main.go

# Distribution
FROM alpine:latest

WORKDIR /dist

COPY --from=builder /builder/bin/chart-viewer .

EXPOSE 9999