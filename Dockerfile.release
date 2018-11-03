FROM alpine:latest
RUN apk update --no-cache && apk  --no-cache upgrade && apk add unzip ca-certificates openssl
WORKDIR /mlxsh
RUN wget https://github.com/ipcjk/mlxsh/releases/download/0.5/release.tar.gz && tar xfz release.tar.gz --strip 1 && rm release.tar.gz && rm /mlxsh/mlxsh.mac /mlxsh/mlxsh.exe && chmod 0755 /mlxsh/mlxsh && chown root:root -R /mlxsh
MAINTAINER Joerg Kost <jk@ip-clear.de>
