FROM golang:latest AS build-env
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod /src/
RUN go mod download
COPY . .
RUN go build -a -o cisco-snmp-pwner -trimpath

FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    && rm -rf /var/cache/*

RUN mkdir -p /app \
    && adduser -D user \
    && chown -R user:user /app

USER user
WORKDIR /app

COPY --from=build-env /src/cisco-snmp-pwner .

ENTRYPOINT [ "/app/cisco-snmp-pwner" ]
