FROM golang:1.16 AS builder

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Download dependencies before build in order to cache them
COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o zeus ./cmd/zeus/main.go

FROM alpine AS app

RUN apk update && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# Copy source files for stacktrace mapping
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

COPY assets ./assets
COPY migrations ./migrations
COPY --from=builder /app/zeus ./

# App
EXPOSE 1111

CMD [ "/app/zeus" ]
