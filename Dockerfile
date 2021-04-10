FROM golang:alpine AS builder

WORKDIR $GOPATH/src/playerdata.co.uk/flake-reporter/

# Create appuser
ENV USER=appuser
ENV UID=1001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/flake-reporter .



FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/flake-reporter /go/bin/flake-reporter

USER appuser:appuser
EXPOSE 9090

ENTRYPOINT ["/go/bin/flake-reporter"]
