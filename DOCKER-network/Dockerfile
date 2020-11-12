FROM ubuntu:16.04

USER root

#Update and install required libs
RUN apt-get update && apt-get install -y \
  apt-utils \
  build-essential \
  git \
  nano \
  tmux \
  vim \
  wget \
  libsnappy-dev

#Setup environment variables
ENV GOVERSION 1.14.7
ENV GOROOT /usr/local/go
ENV GOPATH /home/ubuntu/go
ENV INCL /usr/local/lib
ENV PATH $GOPATH:$GOPATH/bin:$GOROOT/bin:$INCL:$PATH
ENV OLDATA /opt/data/devnet
ENV OLTEST 1
ENV GO111MODULE on

#Create Directories
RUN mkdir -p -- $GOPATH $OLDATA

#Install dependencies for cleveldb
RUN cd /usr/local && wget https://github.com/google/leveldb/archive/v1.20.tar.gz && \
    tar -zxvf v1.20.tar.gz && \
    cd leveldb-1.20/ && \
    make && \
    cp -r out-static/lib* out-shared/lib* /usr/local/lib/ && \
    cd include/ && \
    cp -r leveldb /usr/local/include/ && \
    ldconfig &&\
    cd /usr/local && \
    rm -f v1.20.tar.gz

#Get Golang and install
RUN cd /usr/local && wget https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz

CMD ["/bin/bash"]