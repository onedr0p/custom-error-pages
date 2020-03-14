FROM golang:1.13-alpine as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

RUN apk add --no-cache git ca-certificates

WORKDIR /go/src/github.com/onedr0p/custom-error-pages
COPY . .
RUN apkArch="$(apk --print-arch)"; \
    case "$apkArch" in \
        armv7) export GOARCH='arm' GOARM=7 ;; \
        aarch64) export GOARCH='arm64' ;; \
        x86_64) export GOARCH='amd64' ;; \
    esac; \
    echo "----------------------------------" \
    && echo "apk arch: $(apk --print-arch)" \
    && echo "parsed arch: ${GOARCH}" \
    && echo "----------------------------------" \
    && go build -o custom-error-pages main.go metrics.go \
    && chmod +x custom-error-pages

FROM alpine:3.11

RUN apk add --no-cache ca-certificates tini curl

COPY --from=build /go/src/github.com/onedr0p/custom-error-pages/custom-error-pages /usr/local/bin/custom-error-pages
RUN chmod +x /usr/local/bin/custom-error-pages

ENTRYPOINT ["/sbin/tini", "--", "custom-error-pages"]