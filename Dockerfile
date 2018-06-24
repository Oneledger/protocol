# Use golang Debian-9 image
FROM golang:1.10.3-stretch

USER root

# Get software dependencies
RUN apt-get update && apt-get install -y \
  apt-utils \
  build-essential \
  devscripts \
  git \
  jq \
  libboost-all-dev \
  libevent-dev \
  libssl-dev

# Create oneledger user and its home
RUN rm -rf /var/lib/apt/lists/* \
  && useradd -ms /bin/bash oneledger \
  && mkdir -p /home/oneledger/go/src/github.com/Oneledger/protocol \
  && chown -R oneledger:oneledger /home/oneledger

USER oneledger

# Copy source files over to directory

# Set necessary environment variables for Oneledger/protocol
ENV GOPATH /home/oneledger/go
ENV OLROOT $GOPATH/src/github.com/Oneledger
ENV OLTEST $OLROOT/protocol/node/tests
ENV OLSCRIPT $OLROOT/protocol/node/scripts
ENV OLSETUP $OLROOT/protocol/node/setup
ENV OLDATA $GOPATH/test
# Add user $GOPATH/bin to PATH
ENV PATH $PATH:$GOPATH/bin

COPY --chown=oneledger:oneledger . $OLROOT/protocol

# Set directories for building from src
ENV TENDERMINT_DIR $GOPATH/src/github.com/tendermint/tendermint
ENV BITCOIN_DIR $GOPATH/src/github.com/bitcoin/bitcoin
ENV GETH_DIR $GOPATH/src/github.com/ethereum/go-ethereum

# Set version numbers
ENV TENDERMINT_VERSION v0.18.0
ENV BITCOIN_VERSION v0.16.0

# Install tendermint consensus
RUN mkdir -p $TENDERMINT_DIR \
  && git clone https://github.com/tendermint/tendermint $TENDERMINT_DIR \
  && cd $TENDERMINT_DIR \
  && git checkout tags/$TENDERMINT_VERSION \
  && make get_tools \
  && make get_vendor_deps \
  && make install

# Install geth
RUN mkdir -p $GETH_DIR \
  && git clone https://github.com/ethereum/go-ethereum $GETH_DIR \
  && cd $GETH_DIR \
  && make geth \
  && ln -s $GETH_DIR/build/bin/geth $GOPATH/bin

USER root

# Install bitcoind 0.16
RUN git clone https://github.com/bitcoin/bitcoin $BITCOIN_DIR \
  && cd $BITCOIN_DIR \
  && git checkout $BITCOIN_VERSION

RUN chown -R oneledger:oneledger $BITCOIN_DIR

# Install berkeley libdb 4.8
RUN $BITCOIN_DIR/contrib/install_db4.sh $BITCOIN_DIR

ENV BDB_PREFIX $BITCOIN_DIR/db4

# Compile bitcoin
RUN $BITCOIN_DIR/autogen.sh \
  && $BITCOIN_DIR/configure --without-gui BDB_LIBS="-L${BDB_PREFIX}/lib -ldb_cxx-4.8" BDB_CFLAGS="-I${BDB_PREFIX}/include"

RUN make && make install

USER oneledger

WORKDIR $OLROOT/protocol

