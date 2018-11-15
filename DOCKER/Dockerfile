# Use golang Debian-9 image
FROM golang:1.10.3-stretch

USER root

RUN apt-get update && apt-get install -y \
	curl \
	git \
	jq

RUN rm -rf /var/lib/apt/lists/* \
  && useradd -ms /bin/bash oneledger \
  && mkdir -p /home/oneledger/go/src/github.com/Oneledger/protocol \
  && chown -R oneledger:oneledger /home/oneledger

COPY ./DOCKER/scripts/startnode.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/startnode.sh

USER oneledger

ENV GOPATH /home/oneledger/go
ENV OLROOT $GOPATH/src/github.com/Oneledger
ENV OLTEST $OLROOT/protocol/node/tests
ENV OLSCRIPT $OLROOT/protocol/node/scripts
ENV OLSETUP $OLROOT/protocol/node/setup
ENV OLDATA $GOPATH/test

ENV OL_PORT_RPC 26601
ENV OL_PORT_P2P 26611
ENV OL_PORT_APP 26621
ENV OL_PORT_BTC 18831

EXPOSE ${OL_PORT_RPC}
EXPOSE ${OL_PORT_P2P}
EXPOSE ${OL_PORT_APP}

ENV PATH ${PATH}:${GOPATH}/bin

ENV ID Edwin

ENV TENDERMINT_VERSION v0.23.1
ENV TENDERMINT_DIR $GOPATH/src/github.com/tendermint/tendermint
RUN mkdir -p $TENDERMINT_DIR \
	&& git clone https://github.com/tendermint/tendermint $TENDERMINT_DIR \
	&& cd $TENDERMINT_DIR \
	&& git checkout tags/$TENDERMINT_VERSION \
	&& make get_tools \
	&& make get_vendor_deps \
	&& make install

RUN mkdir -p $OLROOT \
	&& git clone https://github.com/Oneledger/protocol $OLROOT/protocol \
	&& cd $OLROOT/protocol/node \
	&& make update \
	&& make install

WORKDIR /home/oneledger


VOLUME [ "$OLDATA" ]

ENTRYPOINT ["/usr/local/bin/startnode.sh"]

STOPSIGNAL SIGTERM
