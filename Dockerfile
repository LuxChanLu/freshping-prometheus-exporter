FROM golang:1.15.3-alpine3.12 as build

WORKDIR /build/

ENV USER go
ENV UID 10001

RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates && \
    adduser \
        --disabled-password \
        --gecos "" \
        --home "/nonexistent" \
        --shell "/sbin/nologin" \
        --no-create-home \
        --uid "${UID}" \
        "${USER}"

COPY . .

RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /app

FROM scratch

USER go:go

ENTRYPOINT ["/app"]

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /app /app
