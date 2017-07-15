FROM alpine:latest
RUN apk update && apk upgrade
RUN apk add  --no-cache unzip ca-certificates openssl
RUN apk add --no-cache --virtual .build-deps go git libc-dev
WORKDIR /mlxsh
RUN go get github.com/ipcjk/mlxsh
# RUN apk del .build-deps
FROM alpine:latest
WORKDIR /mlxsh
COPY --from=0 /root/go/bin/mlxsh /mlxsh/mlxsh
COPY --from=0 /root/go/src/github.com/ipcjk/mlxsh/mlxsh.yaml  /mlxsh/mlxsh.yaml
MAINTAINER Joerg Kost <jk@ip-clear.de>
