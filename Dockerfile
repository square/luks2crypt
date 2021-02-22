FROM ubuntu:focal

# luks2crypt variables
ARG LUKS2CRYPT_VER='7213ec6894a6f368375a290de81c17f56190c20e'

# golang variables
ARG GO111MODULE='on'
ARG GOLANGVER='1.16'
ARG GOLANGSHA="013a489ebb3e24ef3d915abe5b94c3286c070dfe0818d5bca8108f1d6e8440d2"
ARG GOPATH='/go'
ENV PATH=${PATH}:/usr/local/go/bin:/go/bin

# disable apt interactive prompts
ENV DEBIAN_FRONTEND='noninteractive'

# create golang user
RUN adduser --shell /bin/sh --system --group \
    --gecos 'golang build user' --disabled-password golang

# install golang
RUN apt-get update && apt-get install -y \
      wget \
      ca-certificates \
    && wget --no-verbose "https://dl.google.com/go/go${GOLANGVER}.linux-amd64.tar.gz" \
    && echo "${GOLANGSHA} *go${GOLANGVER}.linux-amd64.tar.gz" > go${GOLANGVER}.linux-amd64.tar.gz.shasum \
    && sha256sum -c go${GOLANGVER}.linux-amd64.tar.gz.shasum \
    && tar -C /usr/local -xzf go${GOLANGVER}.linux-amd64.tar.gz \
    && mkdir -p /go \
    && chown -R golang:golang /go \
    && apt-get remove -y wget \
    && rm -f go${GOLANGVER}.linux-amd64.tar.gz \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# install luks2crypt build dependencies
RUN apt-get update && apt-get install -y \
      build-essential \
      pkg-config \
      git \
      libcryptsetup-dev \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# set build user to golang
USER golang

# create dir structure for luks2crypt build
RUN mkdir -p /go/src/github.com/square/luks2crypt
WORKDIR /go/src/github.com/square/luks2crypt

# copy in files
COPY . .

# go build and install luks2crypt
RUN make install

# run and print the version of luks2crypt
ENTRYPOINT [ "luks2crypt" ]
CMD [ "version" ]
