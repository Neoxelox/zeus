FROM golang:1.16 AS builder

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Download dependencies before build in order to cache them
COPY go.mod go.sum /app/

RUN go mod download

COPY . /app

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o zeus ./cmd/zeus/main.go

FROM alpine AS app

RUN apk update && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# Copy source files for stacktrace mapping
COPY cmd /app/cmd
COPY internal /app/internal
COPY pkg /app/pkg

COPY assets /assets
COPY --from=builder /app/zeus /app

# App
EXPOSE 1111

CMD [ "/app/zeus" ]