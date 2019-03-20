FROM \
    karalabe/xgo-latest

MAINTAINER \
    Vincent <vincent at renproject.io>

# only add 64-bit architectures since we're only building for amd64
# and also install curl
RUN \
    dpkg --add-architecture amd64 && \
    apt-get update && \
    apt-get install -y --no-install-recommends curl

