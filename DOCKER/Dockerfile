FROM golang:1.10.3-stretch

USER root

RUN apt-get update && apt-get install -y \
  apt-utils \
  build-essential \
  git \
  nano \
  tmux \
  vim

RUN rm -rf /var/lib/apt/lists/* \
  && useradd -ms /bin/bash oneledger \
  && mkdir -p /home/oneledger/go/src/github.com/Oneledger/protocol \
  && mkdir -p /home/oneledger/testnet \
  && chown -R oneledger:oneledger /home/oneledger

USER oneledger

ENV GOPATH /home/oneledger/go
ENV OLROOT $GOPATH/src/github.com/Oneledger
ENV OLTEST $OLROOT/protocol/node/tests
ENV OLSCRIPT $OLROOT/protocol/node/scripts
ENV OLSETUP $OLROOT/protocol/node/setup
ENV OLDATA /home/oneledger/.olfullnode
ENV OLVERSION v0.8.1
ENV PATH $PATH:$GOPATH/bin

USER oneledger

RUN git clone https://github.com/Oneledger/protocol $OLROOT/protocol
RUN chown -R oneledger:oneledger $OLROOT/protocol

WORKDIR $OLROOT/protocol

RUN git checkout $OLVERSION

USER root
RUN cp $OLROOT/protocol/DOCKER/chronos/bootstrapNode /usr/local/bin \
	&& cp $OLROOT/protocol/DOCKER/chronos/startNode /usr/local/bin \
	&& cp $OLROOT/protocol/DOCKER/chronos/stopNode /usr/local/bin \
	&& cp $OLROOT/protocol/DOCKER/chronos/cleanNode /usr/local/bin
USER oneledger

WORKDIR $OLROOT/protocol/node

RUN make tools update install

VOLUME [ "$OLDATA" ]
