FROM golang:1.17-alpine AS builder
RUN apk update && apk add --no-cache make git
WORKDIR /go/src/gitlab.com/rarimo/bdjuno
COPY . ./
RUN go mod download
RUN make build

FROM alpine:latest
WORKDIR /bdjuno
COPY --from=builder /go/src/gitlab.com/rarimo/bdjuno/build/bdjuno /usr/bin/bdjuno
CMD [ "bdjuno" ]
