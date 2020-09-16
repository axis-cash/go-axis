#!/bin/sh



LOCAL_PATH=$(cd `dirname $0`; pwd)
echo "LOCAL_PATH=$LOCAL_PATH"
AXIS_PATH="${LOCAL_PATH%}"
echo "AXIS_PATH=$AXIS_PATH"
CZERO_PATH="${AXIS_PATH%/*}/go-axis-import"
echo "CZERO_PATH=$CZERO_PATH"

echo "update go-axis-import"
cd $CZERO_PATH
git fetch&&git rebase

echo "update go-axis"
cd $AXIS_PATH
git fetch&&git rebase
make clean
BUILD_PATH="${AXIS_PATH%}/build"

os="all"
version="v0.3.1-beta.rc.5"
while getopts ":o:v:" opt
do
    case $opt in
        o)
        os=$OPTARG
        ;;
        v)
        version=$OPTARG
        ;;
        ?)
        echo "unkonw param"
        exit 1;;
    esac
done

if [ "$os" = "all" ]; then
    os_version=("linux-amd64-v3" "linux-amd64-v4" "darwin-amd64" "windows-amd64")
else
    os_version[0]="$os"
fi

for os in ${os_version[@]}
    do
      echo "make gaxistx-${os}"
      make "gaxistx-"${os}
      rm -rf $BUILD_PATH/gaxistxpkg/bin
      rm -rf $BUILD_PATH/gaxistxpkg/czero
      mkdir -p $BUILD_PATH/gaxistxpkg/bin
      mkdir -p $BUILD_PATH/gaxistxpkg/czero/data/
      mkdir -p $BUILD_PATH/gaxistxpkg/czero/include/
      mkdir -p $BUILD_PATH/gaxistxpkg/czero/lib/
      cp -rf $CZERO_PATH/czero/data/* $AXIS_PATH/build/gaxistxpkg/czero/data/
      cp -rf $CZERO_PATH/czero/include/* $AXIS_PATH/build/gaxistxpkg/czero/include/
      if [ $os == "windows-amd64" ];then
        mv $BUILD_PATH/bin/gaxistx*.exe $BUILD_PATH/gaxistxpkg/bin/tx.exe
        cp -rf  $CZERO_PATH/czero/lib_WINDOWS_AMD64/* $AXIS_PATH/build/gaxistxpkg/czero/lib/
      elif [ $os == "linux-amd64-v3" ];then
        mv $BUILD_PATH/bin/gaxistx-v3* $BUILD_PATH/gaxistxpkg/bin/tx
        cp -rf  $CZERO_PATH/czero/lib_LINUX_AMD64_V3/* $AXIS_PATH/build/gaxistxpkg/czero/lib/
      elif [ $os == "linux-amd64-v4" ];then
        mv $BUILD_PATH/bin/gaxistx-v4* $BUILD_PATH/gaxistxpkg/bin/tx
        cp -rf  $CZERO_PATH/czero/lib_LINUX_AMD64_V4/* $AXIS_PATH/build/gaxistxpkg/czero/lib/
      else
        mv $BUILD_PATH/bin/gaxistx-darwin* $BUILD_PATH/gaxistxpkg/bin/tx
        cp -rf  $CZERO_PATH/czero/lib_DARWIN_AMD64/* $AXIS_PATH/build/gaxistxpkg/czero/lib/
      fi
      cd $BUILD_PATH

      if [ $os == "windows-amd64" ];then
        rm -rf ./gaxistx-*-$os.zip
        zip -r gaxistx-$version-$os.zip gaxistxpkg/*
      else
         rm -rf ./gaxistx-*-$os.tar.gz
         tar czvf gaxistx-$version-$os.tar.gz gaxistxpkg/*
      fi

      cd $LOCAL_PATH

    done
rm -rf $BUILD_PATH/gaxistxpkg/bin
rm -rf $BUILD_PATH/gaxistxpkg/czero

