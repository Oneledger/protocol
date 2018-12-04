#!/bin/bash

# Install dependencies using Homebrew on Mac
install_dependencies_on_mac () {
	brew install autoconf
	brew install automake
	brew install libtool
	brew install boost
	brew install miniupnpc
	brew install openssl
	brew install pkg-config
	brew install protobuf
	brew install libevent
	# Install berkeley-db
	brew install berkeley-db4
	# Install Qt5
	#brew install qt5
}

# Building bitcoin on Mac
build_bitcoin_on_mac () {
	# Clone source code
	git clone https://github.com/bitcoin/bitcoin.git

	cd bitcoin
	git checkout v0.16.0

	# Build
	./autogen.sh
	# Disable the GUI
	./configure --without-gui --disable-tests
	make -j$(nproc)
	make install
}

# Move bitcoind on Mac
move_bitcoind_on_mac () {
	# Clone source code
	git clone https://github.com/bitcoin/bitcoin.git
	cd bitcoin
	# Build
	./autogen.sh
	# Disable the GUI
	./configure --without-gui --disable-tests
	make -j$(nproc)
}

install_dependencies_on_linux () {
	sudo apt-get install -y build-essential libtool autotools-dev automake pkg-config bsdmainutils python3
	sudo apt-get install -y libssl-dev libevent-dev libboost-all-dev libboost-system-dev libboost-filesystem-dev libboost-chrono-dev libboost-test-dev libboost-thread-dev

}

build_bitcoin_on_linux () {
	# Clone source code
	git clone https://github.com/bitcoin/bitcoin.git

	cd bitcoin
	git checkout v0.16.0
    ./contrib/install_db4.sh $BITCOIN_DIR
    BDB_PREFIX=$BITCOIN_DIR/db4
	# Build
	./autogen.sh
	# Disable the GUI
	./configure --without-gui --disable-tests BDB_LIBS="-L${BDB_PREFIX}/lib -ldb_cxx-4.8" BDB_CFLAGS="-I${BDB_PREFIX}/include"
	make -j$(nproc)
	sudo make install
}

current_uname=$(uname);
case "$current_uname" in
    (*Darwin*)
	current_platform='MacOS'
	which -s brew
	if [[ $? != 0 ]] ; then
		# Install Homebrew
		/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
	else
		brew update
	fi

	# Test if bitcoind installed
	which -s bitcoind
	if [[ $? != 0 ]] ; then
		echo 'Install dependencies'
		install_dependencies_on_mac
		if [[ $? == 0 ]] ; then
			echo 'Building bitcoin'
			build_bitcoin_on_mac
		fi
	fi
	;;
    (*Linux*)
	current_platform='Linux'
	sudo apt-get update
	sudo apt-get install -y software-properties-common

	# Test if bitcoind installed
	which bitcoind
	if [[ $? != 0 ]] ; then
#		# Add bitcoin official PPA repo
#		sudo add-apt-repository -y ppa:bitcoin/bitcoin
#		sudo apt-get update
#		# Install bitcoind from PPA
#		sudo apt-get install bitcoind
        echo "Installing Bitcoin"
        TMPFOLD=/tmp
        BITCOIN_DIR=$TMPFOLD/bitcoin
        pushd $TMPFOLD > /dev/null
        echo 'Install dependencies'
		install_dependencies_on_linux
		if [[ $? == 0 ]] ; then
		    echo 'Building bitcoin'
		    build_bitcoin_on_linux
		fi
		popd > /dev/null
	fi
	;;
    (*) echo 'error: unsupported platform.'; exit 2; ;;
esac;

bitcoind -version
if [[ $? == 0 ]] ; then
	echo 'bitcoind was installed successfully'
fi
