# Step 1: Modules caching
FROM golang:1.26-alpine AS modules

COPY go.mod go.sum /modules/

ARG GITHUB_TOKEN

WORKDIR /modules

RUN apk add --no-cache git && \
    git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/" && \
    GOPRIVATE=github.com/Muxly-Corp go mod download

# Step 2: Builder
FROM golang:1.26-alpine AS builder

COPY --from=modules /go/pkg /go/pkg

COPY . /build/

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/muxly-msg-subscriber ./cmd/app

# Step 3: Production
FROM scratch AS prod

COPY --from=builder /bin/muxly-msg-subscriber /muxly-msg-subscriber
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/muxly-msg-subscriber"]

# Step 4: Development
FROM alpine AS dev

COPY --from=builder /bin/muxly-msg-subscriber /bin/muxly-msg-subscriber
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY dev_entrypoint.sh /dev_entrypoint.sh

RUN chmod +x /dev_entrypoint.sh

ENTRYPOINT ["/dev_entrypoint.sh"]
