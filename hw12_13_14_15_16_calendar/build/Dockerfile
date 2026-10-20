FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG LDFLAGS=""
RUN CGO_ENABLED=0 GOOS=linux go build \
        -trimpath \
        -ldflags "$LDFLAGS" \
        -o /out/calendar ./cmd/calendar

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 app

COPY --from=build /out/calendar /usr/local/bin/calendar
COPY configs/config.yaml /etc/calendar/config.yaml

USER app

ENTRYPOINT ["/usr/local/bin/calendar"]
CMD ["-config", "/etc/calendar/config.yaml"]
