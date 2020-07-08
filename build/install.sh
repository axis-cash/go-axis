#! /bin/sh

_GOPATH=`cd ../../../../../;pwd`

export GOPATH=$_GOPATH
echo $GOPATH

go install -v ../cmd/gaxis
go install -v ../cmd/axiskey
