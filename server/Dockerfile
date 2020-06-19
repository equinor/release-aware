FROM golang:1.12.5

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# install required debian packages
# add any package that is required after `build-essential`, end the line with \
RUN apt-get update && apt-get install -y \
    build-essential \
&& rm -rf /var/lib/apt/lists/*

# install glide,  gin and delve
RUN go get github.com/Masterminds/glide
RUN go get github.com/codegangsta/gin
RUN go get github.com/go-delve/delve/cmd/dlv

# setup the working directory
WORKDIR /go/src/app
ADD glide.yaml glide.yaml
ADD glide.lock glide.lock
RUN mkdir /scripts
ADD run-local-server.sh /scripts/run-local-server.sh
RUN chmod +x /scripts/run-local-server.sh

# install dependencies
RUN glide install --skip-test

# add source code
ADD src src

# build the source
RUN go build src/main.go

# Expose 2345 for Delve debugger
EXPOSE 2345