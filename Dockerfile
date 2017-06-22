FROM alpine:3.5
RUN apk update
RUN apk add unzip
RUN apk add ca-certificates
RUN apk add openssl
WORKDIR /brocadecli
RUN wget -O /brocadecli/brocadecli.linux https://github.com/ipcjk/brocadecli/raw/master/bin/brocadecli.linux
RUN chmod 0755 /brocadecli/brocadecli.linux
MAINTAINER Joerg Kost <jk@ip-clear.de>
