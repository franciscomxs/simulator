FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

ENV GOFLAGS="-mod=mod"
ENV GONOSUMDB="*"

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/web ./cmd/web && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/lambda ./cmd/lambda

FROM alpine:3.19

RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=builder /bin/web /app/web

EXPOSE 8080
CMD ["/app/web"]
