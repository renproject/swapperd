# Custom Dockerfile that uses xgo
# Based off:
# https://github.com/billziss-gh/cgofuse/blob/9b5a7c093a2b5da9dc74494e4c7714af8c82de93/Dockerfile

FROM \
    karalabe/xgo-latest

MAINTAINER \
    Vincent <vincent at renproject.io>

# only add 64-bit architectures since we're only building for amd64
# and also install curl and the latest version of go
RUN \
    dpkg --add-architecture amd64 && \
    apt-get install software-properties-common && \
    add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update && \
    apt-get install -y --no-install-recommends curl golang-go

