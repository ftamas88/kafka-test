FROM golang:1.16.2-alpine3.13 as build

COPY . /source
WORKDIR /source

ARG GITHUB_OAUTH_TOKEN
RUN git config --global url."https://$GITHUB_OAUTH_TOKEN:@github.com/".insteadOf "https://github.com/"

RUN go build -mod=readonly -ldflags "-X github.com/ftamas88/kafka-test/internal/app.version=`git describe --tags --dirty` -X github.com/thisisbud/ftamas88/kafka-test/internal/app.commitHash=`git rev-parse HEAD`" -o bin/kafka-test -v cmd/app/*.go
FROM alpine:3.13

COPY --from=build /source/bin/kafka-test /source/bin/kafka-test
COPY --from=build /source/configs /source/configs

WORKDIR /source
ENTRYPOINT ["./bin/kafka-test"]
