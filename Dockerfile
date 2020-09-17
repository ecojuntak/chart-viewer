# Builder
FROM golang:1.14 as builder
LABEL stage=builder

WORKDIR /builder

COPY . .

RUN CGO_ENABLED=0 go build -o ./bin/app .

# Distribution
FROM alpine:latest

WORKDIR /dist

COPY --from=builder /builder/bin/app .

EXPOSE 9999
CMD /dist/app