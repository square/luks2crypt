FROM golang:1.11

ARG LUKS2CRYPT_VER='7213ec6894a6f368375a290de81c17f56190c20e'
ARG GO111MODULE='on'

RUN adduser --shell /bin/sh --system --group \
    --gecos 'golang build user' --disabled-password golang

RUN apt-get update && apt-get install -y \
      libcryptsetup-dev \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

USER golang

RUN mkdir -p /go/src/github.com/square/luks2crypt
WORKDIR /go/src/github.com/square/luks2crypt

COPY . .

RUN go install -ldflags "-X main.VERSION=${LUKS2CRYPT_VER}" -v ./...

ENTRYPOINT [ "luks2crypt" ]
CMD [ "version" ]
