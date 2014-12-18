#!/bin/bash -e
if [[ ! -f /etc/debian_version ]];then
    echo Only debian based supported
    exit 1
fi

if [[ ! -d ~/.gvm ]];then
    apt-get install -qq -yy git mercurial  bison make curl
    wget -q https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer -O /tmp/gvm-installer
    bash /tmp/gvm-installer
fi

source ~/.gvm/scripts/gvm
gvm install go1.3
gvm use go1.3 --global

go get github.com/mattn/gom
