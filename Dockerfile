FROM alpine:3.5
RUN apk update
RUN apk add unzip
RUN apk add ca-certificates
RUN apk add openssl
WORKDIR /mlxsh
RUN wget -O /mlxsh/mlxsh https://github.com/ipcjk/mlxsh/raw/master/bin/mlxsh
RUN chmod 0755 /mlxsh/mlxsh
MAINTAINER Joerg Kost <jk@ip-clear.de>
