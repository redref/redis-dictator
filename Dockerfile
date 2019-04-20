FROM golang:1.11-alpine as builder

WORKDIR $GOPATH/src/github.com/blablacar/redis-dictator
COPY . ./
RUN \
  go build ./... && \
  mv dictator /

FROM alpine

COPY --from=builder /dictator /

CMD ["/dictator", "/etc/dictator/dictator.json"]
