# Backend builder
FROM golang:1.16 as backend-builder
LABEL stage=builder
WORKDIR /builder

COPY . .

RUN make build-backend

# Frontend builder
FROM node:14-buster as frontend-builder
LABEL stage=builder
WORKDIR /builder

COPY . .

RUN make build-frontend

# Distribution
FROM alpine:latest
WORKDIR /app

COPY --from=backend-builder /builder/bin/chart-viewer .
COPY --from=backend-builder /builder/seed.json ./seed.json
COPY --from=backend-builder /builder/api_versions.json ./api_versions.json
COPY --from=frontend-builder /builder/ui/dist ./ui/dist
