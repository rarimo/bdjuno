FROM golang:1.18-alpine as buildbase

WORKDIR /go/src/gitlab.com/rarimo/bdjuno
COPY vendor .
COPY . .
ENV GO111MODULE="on"
ENV CGO_ENABLED=1
ENV GOOS="linux"
RUN go build -o /usr/local/bin/bdjuno gitlab.com/rarimo/bdjuno

###

FROM alpine:3.9
COPY --from=buildbase /usr/local/bin/bdjuno /usr/local/bin/bdjuno

ENTRYPOINT ["bdjuno"]
