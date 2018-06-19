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
	./configure --without-gui
	make
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
	./configure --without-gui
	make
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
	apt-get update
	apt-get install software-properties-common

	# Test if bitcoind installed
	which bitcoind
	if [[ $? != 0 ]] ; then
		# Add bitcoin official PPA repo
		add-apt-repository -y ppa:bitcoin/bitcoin
		apt-get update
		# Install bitcoind from PPA
		apt-get install bitcoind
	fi
	;;
    (*) echo 'error: unsupported platform.'; exit 2; ;;
esac;

bitcoind -version
if [[ $? == 0 ]] ; then
	echo 'bitcoind was installed successfully'
fi
