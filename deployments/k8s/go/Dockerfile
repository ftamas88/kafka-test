FROM golang:1.16.2-alpine3.13 as build

COPY . /source
WORKDIR /source

# ARG GITHUB_OAUTH_TOKEN
# RUN git config --global url."https://$GITHUB_OAUTH_TOKEN:@github.com/".insteadOf "https://github.com/"

RUN sed -i 's/http\:\/\/dl-cdn.alpinelinux.org/http\:\/\/dl-4.alpinelinux.org/g' /etc/apk/repositories

RUN addgroup -g 3000 -S tamas
RUN adduser -u 3000 -D -S tamas -G tamas

RUN apk add git gcc g++ libc-dev make yarn
RUN wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest
RUN go get golang.org/x/tools/cmd/goimports
RUN yarn global add speccy

# RUN go build -mod=readonly -ldflags "-X github.com/ftamas88/kafka-test/internal/app.version=`git describe --tags --dirty` -X github.com/thisisbud/ftamas88/kafka-test/internal/app.commitHash=`git rev-parse HEAD`" -o bin/kafka-test -v cmd/app/*.go
RUN go build -mod=readonly -o bin/kafka-test -v cmd/app/app.go

FROM alpine:3.13

COPY --from=build /source/bin/kafka-test /source/bin/kafka-test
COPY --from=build /source/configs /source/configs

WORKDIR /source
ENTRYPOINT ["./bin/kafka-test"]
