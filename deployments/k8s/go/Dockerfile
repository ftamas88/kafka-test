FROM golang:1.15.10-alpine3.13 as build

# ARG GITHUB_OAUTH_TOKEN
# RUN git config --global url."https://$GITHUB_OAUTH_TOKEN:@github.com/".insteadOf "https://github.com/"

RUN sed -i 's/http\:\/\/dl-cdn.alpinelinux.org/http\:\/\/dl-4.alpinelinux.org/g' /etc/apk/repositories

#ENV CGO_ENABLED=0
#ENV go enw -w CGO_ENABLED=0
RUN go env -w GOPRIVATE=github.com/ftamas88/*
RUN go env -w  GO111MODULE=auto

RUN addgroup -g 3000 -S tamas
RUN adduser -u 3000 -D -S tamas -G tamas

RUN apk add git gcc g++ # libc-dev make yarn
RUN wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest
RUN go get golang.org/x/tools/cmd/goimports
# RUN yarn global add speccy

COPY . /source
WORKDIR /source

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o ./bin/app -v ./cmd/app/app.go

FROM alpine:3.13

COPY --from=build /source/bin/app /app

WORKDIR /source
ENTRYPOINT ["/app"]
