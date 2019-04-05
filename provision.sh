#!/bin/bash
set -o errexit;
set -o nounset;

export GOPACK=go1.12.1.linux-386.tar.gz;
export TARGET_PACK=${HOME}/golang;
export TARGET_GO=/usr/local/go;
export ORIGIN=https://storage.googleapis.com/golang/${GOPACK};

echo 'Installing GCC...';
sudo apt-get update;
sudo apt-get install build-essential -y;

echo 'Creating golang pack directory...';
mkdir -p ${TARGET_PACK};

echo 'Downloading golang pack...';
[ -f ${TARGET_PACK}/${GOPACK} ] || wget -nv --progress=dot:giga -O ${TARGET_PACK}/${GOPACK} ${ORIGIN};

echo 'Install golang binaries ... target: /usr/local';
sudo tar zxf ${TARGET_PACK}/${GOPACK} -C /usr/local/;
sudo chown -R ${USER} ${TARGET_GO};

echo 'Removing golang pack...';
rm -rf ${TARGET_PACK};

echo 'Installing packages...';
${TARGET_GO}/bin/go get -u -v golang.org/x/lint/golint;
${TARGET_GO}/bin/go get -u -v golang.org/x/tools/cmd/godoc;

echo 'Creating directories...';
chown -R ${USER}:${USER} /usr/local/go;
mkdir -p ${HOME}/go/{pkg,bin};
mkdir -p ${HOME}/go/src/github.com/video-sharing-platform-thingy;
ln -s /vagrant ${HOME}/go/src/github.com/video-sharing-platform-thingy/backend;

echo 'Writing a config...';
cat >> ${HOME}/.profile <<- 'EOM'
[ -d /usr/local/go ] && export GOROOT=/usr/local/go;
[ -d ${HOME}/go ] && export GOPATH=${HOME}/go;
[ -d ${GOPATH}/bin ] && export GOBIN=${GOPATH}/bin;
[ -d ${GOROOT}/bin ] && export LPATH=${GOROOT}/bin;
PATH=${PATH}:${LPATH}:${GOBIN};

function run() {
  cd ${HOME}/go/src/github.com/video-sharing-platform-thingy/backend;
  echo 'Downloading dependencies...';
  go get;
  echo 'Starting...';
  go build;
  ./backend;
  rm backend;
}

echo 'Welcome to the VSPT backend Vagrant instance. Run `run` to start the backend.';
EOM

echo 'Done! Run `vagrant ssh` to get started.';