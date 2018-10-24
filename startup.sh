#! /bin/bash
# Description: Setup go environment, pull the go-dumper repo, run the process and finally shutdown the instance 

# Install mysql
sudo apt install mysql-client-5.7 

# Download golang and install
curl -O https://storage.googleapis.com/golang/go1.10.3.linux-amd64.tar.gz
tar xvf go1.10.3.linux-amd64.tar.gz
sudo chown -R root:root ./go
sudo mv go /usr/local

## Set go environment
export GOPATH=$HOME/work
export GOBIN=$GOPATH/bin
export PATH=$PATH:/usr/local/go/bin:$GOBIN

## Crerate go work directories
mkdir $HOME/work
mkdir $HOME/work/bin
mkdir $HOME/work/pkg
mkdir $HOME/work/src

dstDir="$HOME/work/src/github.com/go-dumper"

mkdir -p $dstDir
cd $dstDir

## Clone the repo
git clone https://DatpTK4yc_dyUU5QzKLh@github.com/go-dumper.git

## Start the backup process
source run.sh prod

## Shutdown instance. Note: this just shuts down the instance-not delete it.
sudo shutdown -h now