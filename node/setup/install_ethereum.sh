#!/bin/bash

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

	# Test if geth installed
	which -s geth
	if [[ $? != 0 ]] ; then
		# Add the tap and install geth
		brew tap ethereum/ethereum
		brew install --verbose --debug ethereum
	fi
	;;
    (*Linux*)
	current_platform='Linux'
	sudo apt-get update
	sudo apt-get install software-properties-common

	# Test if geth installed
	which geth
	if [[ $? != 0 ]] ; then
		# Install geth from PPA
		sudo add-apt-repository -y ppa:ethereum/ethereum
		sudo apt-get update
		sudo apt-get install ethereum
	fi
	;;
    (*) echo 'error: unsupported platform.'; exit 2; ;;
esac;

geth version
echo 'geth was installed successfully'
