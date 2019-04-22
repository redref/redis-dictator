FROM golang:1.11-alpine as builder

ARG VERSION
ARG BUILD_DATE
ARG GIT_REVISION
WORKDIR $GOPATH/src/github.com/Junonogis/redis-dictator
COPY . ./
RUN \
  go build -ldflags "-X main.BuildTime=${BUILD_DATE} -X main.Version=${VERSION} -X main.GitRevision=${GIT_REVISION}" ./... && \
  mv dictator /

FROM alpine

COPY --from=builder /dictator /
COPY example/dictator.json /etc/dictator/dictator.json

CMD ["/dictator", "--config", "/etc/dictator/dictator.json"]
