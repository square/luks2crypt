FROM golang:1.9

ARG LUKS2CRYPT_VER='7213ec6894a6f368375a290de81c17f56190c20e'
ARG DEP_VERSION='v0.3.2'

RUN adduser --shell /bin/sh --no-create-home --system --group \
    --gecos 'golang build user' --disabled-password golang

RUN apt-get update && apt-get install -y \
      libcryptsetup-dev \
      curl \
    && curl -fsSL https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64 -o $GOPATH/bin/dep \
    && chmod +x $GOPATH/bin/dep \
    && apt-get remove -y curl \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

USER golang

RUN mkdir -p /go/src/github.com/square/luks2crypt
WORKDIR /go/src/github.com/square/luks2crypt

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .

# "go install -v ./..."
RUN go-wrapper install -ldflags "-X main.VERSION=${LUKS2CRYPT_VER}"

# ["app"]
ENTRYPOINT [ "go-wrapper", "run" ]
