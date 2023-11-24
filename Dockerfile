FROM golang:1.20-alpine as buildbase

WORKDIR /go/src/github.com/rarimo/bdjuno
RUN apk add build-base
COPY . .
ENV GO111MODULE="on"
ENV CGO_ENABLED=1
ENV GOOS="linux"
RUN go mod tidy
RUN go mod vendor
RUN go build -o /usr/local/bin/bdjuno github.com/rarimo/bdjuno/cmd/bdjuno

###

FROM alpine:3.9
COPY --from=buildbase /usr/local/bin/bdjuno /usr/local/bin/bdjuno

ENTRYPOINT ["bdjuno"]
