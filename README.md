restartd
=======

restartd allow systemd service control through unix socket.

Install
=======

Fetch sources

    git clone git@github.com:bearstech/restartd.git $GOPATH/src/restartd

Get deps

    cd $GOPATH/src/restartd && make get

Build

    cd $GOPATH/src/restard && make

Install

    cd $GOPATH/src/restartd && make install ROOTFS=/path/to/rootfs

