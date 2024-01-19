# Stage: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0

COPY . .

RUN go build -trimpath -o app .

# Stage: Final
FROM ubuntu:22.10

WORKDIR /app
COPY --from=builder /app/app /app/app

ENTRYPOINT [ "/app/app" ]
