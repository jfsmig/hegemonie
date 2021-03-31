#------------------------------------------------------------------------------
# Install the system dependencies

FROM golang:1.16-buster AS dependencies

LABEL maintainer="Jean-Francois SMIGIELSKI <jf.smigielski@gmail.com>"

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPATH=/gopath

USER 0
WORKDIR /

RUN set -ex \
&& apt-get update -y \
&& apt-get install -y --no-install-recommends \
  make \
  protobuf-compiler \
  librocksdb-dev \
  librocksdb5.17 \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*

RUN set -ex \
&& go get -u -t github.com/golang/protobuf@v1.4.3 \
&& go get -u -t google.golang.org/grpc@v1.35.0 \
&& go get \
  google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0 \
  google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0

RUN echo 'PATH=/gopath/bin:/go/bin:/usr/local/go/bin:$PATH\nexport PATH' >> /etc/profile



#------------------------------------------------------------------------------
# Install the golang dependencies on top the the system dependencies

FROM dependencies AS builder

USER 0
WORKDIR /gopath/src/github.com/hegemonie-rpg/engine

# Build & Install the code, then extract all the system deps. Inspired by:
# https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i
COPY go.sum go.mod /gopath/src/github.com/hegemonie-rpg/engine/
RUN go mod download

COPY Makefile LICENSE AUTHORS.md /gopath/src/github.com/hegemonie-rpg/engine/
COPY pkg /gopath/src/github.com/hegemonie-rpg/engine/pkg
COPY api /gopath/src/github.com/hegemonie-rpg/engine/api
RUN set -x \
&& echo $PATH \
&& . /etc/profile \
&& echo $PATH \
&& make clean \
&& make prepare \
&& make \
&& mkdir /dist \
&& cp -p -v /gopath/bin/hege /dist/ \
&& cp -p -v /gopath/bin/hege /usr/bin/

WORKDIR /dist

RUN set -ex \
&& mkdir -p /dist/lib64 \
&& ldd /dist/hege | tr -s '[:blank:]' '\n' | grep '^/' | sort | uniq \
 | xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;' \
&& cp /lib64/ld-linux-x86-64.so.2 /dist/lib64/



#------------------------------------------------------------------------------
# Create the minimal runtime image

FROM scratch as runtime

USER 0
COPY --chown=0:0 --from=builder /dist /

EXPOSE 6000
EXPOSE 6001

USER 0
WORKDIR /
ENTRYPOINT ["/hege"]

