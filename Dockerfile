FROM alpine:3.5
RUN apk update
RUN apk add unzip
RUN apk add ca-certificates
RUN apk add openssl
WORKDIR /mlxsh
RUN wget https://github.com/ipcjk/mlxsh/releases/download/0.1/release.tar.gz
RUN tar xfz release.tar.gz --strip 1
RUN rm release.tar.gz
RUN chmod 0755 /mlxsh/mlxsh.mac && chmod 0755 /mlxsh/mlxsh
RUN chown root:root -R /mlxsh
MAINTAINER Joerg Kost <jk@ip-clear.de>
