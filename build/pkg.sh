#!/bin/sh

LOCAL_PATH=$(cd `dirname $0`; pwd)
AXIS_PATH="${LOCAL_PATH%/*}"
CZERO_PATH="${AXIS_PATH%/*}/go-axis-import"

echo "update go-axis-import"
cd $CZERO_PATH
git fetch&&git rebase

echo "update go-axis"
cd $AXIS_PATH
git fetch&&git rebase
make clean all

rm -rf $LOCAL_PATH/gaxispkg/bin
rm -rf $LOCAL_PATH/gaxispkg/czero
mkdir -p $LOCAL_PATH/gaxispkg/czero/data/
mkdir -p $LOCAL_PATH/gaxispkg/czero/include/
mkdir -p $LOCAL_PATH/gaxispkg/czero/lib/
cp -rf $LOCAL_PATH/bin $LOCAL_PATH/gaxispkg
cp -rf $CZERO_PATH/czero/data/* $AXIS_PATH/build/gaxispkg/czero/data/
cp -rf $CZERO_PATH/czero/include/* $AXIS_PATH/build/gaxispkg/czero/include/

function sysname() {

    SYSTEM=`uname -s |cut -f1 -d_`

    if [ "Darwin" == "$SYSTEM" ]
    then
        echo "Darwin"

    elif [ "Linux" == "$SYSTEM" ]
    then
        name=`uname  -r |cut -f1 -d.`
        echo Linux-V"$name"
    else
        echo "$SYSTEM"
    fi
}

SNAME=`sysname`

if [ "Darwin" == "$SNAME" ]
then
    echo $SNAME
    cp $CZERO_PATH/czero/lib_DARWIN_AMD64/* $AXIS_PATH/build/gaxispkg/czero/lib/
elif [ "Linux-V3" == "$SNAME" ]
then
    echo $SNAME
    cp $CZERO_PATH/czero/lib_LINUX_AMD64_V3/* $AXIS_PATH/build/gaxispkg/czero/lib/
elif [ "Linux-V4" == "$SNAME" ]
then
    echo $SNAME
    cp $CZERO_PATH/czero/lib_LINUX_AMD64_V4/* $AXIS_PATH/build/gaxispkg/czero/lib/
fi

cd $LOCAL_PATH
if [ -f ./gaxispkg_*.tar.gz ]; then
	rm ./gaxispkg_*.tar.gz
fi
tar czvf gaxispkg_$SNAME.tar.gz gaxispkg/*
