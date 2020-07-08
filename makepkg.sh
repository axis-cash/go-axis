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
      echo "make gaxis-${os}"
      make "gaxis-"${os}
      rm -rf $BUILD_PATH/geropkg/bin
      rm -rf $BUILD_PATH/geropkg/czero
      mkdir -p $BUILD_PATH/geropkg/bin
      mkdir -p $BUILD_PATH/geropkg/czero/data/
      mkdir -p $BUILD_PATH/geropkg/czero/include/
      mkdir -p $BUILD_PATH/geropkg/czero/lib/
      cp -rf $CZERO_PATH/czero/data/* $AXIS_PATH/build/geropkg/czero/data/
      cp -rf $CZERO_PATH/czero/include/* $AXIS_PATH/build/geropkg/czero/include/
      if [ $os == "windows-amd64" ];then
        mv $BUILD_PATH/bin/gaxis*.exe $BUILD_PATH/geropkg/bin/gaxis.exe
        cp -rf  $CZERO_PATH/czero/lib_WINDOWS_AMD64/* $AXIS_PATH/build/geropkg/czero/lib/
      elif [ $os == "linux-amd64-v3" ];then
        mv $BUILD_PATH/bin/bootnode-v3*  $BUILD_PATH/geropkg/bin/bootnode
        mv $BUILD_PATH/bin/gaxis-v3* $BUILD_PATH/geropkg/bin/gaxis
        cp -rf  $CZERO_PATH/czero/lib_LINUX_AMD64_V3/* $AXIS_PATH/build/geropkg/czero/lib/
      elif [ $os == "linux-amd64-v4" ];then
        mv $BUILD_PATH/bin/bootnode-v4*  $BUILD_PATH/geropkg/bin/bootnode
        mv $BUILD_PATH/bin/gaxis-v4* $BUILD_PATH/geropkg/bin/gaxis
        cp -rf  $CZERO_PATH/czero/lib_LINUX_AMD64_V4/* $AXIS_PATH/build/geropkg/czero/lib/
      else
        mv $BUILD_PATH/bin/gaxis-darwin* $BUILD_PATH/geropkg/bin/gaxis
        cp -rf  $CZERO_PATH/czero/lib_DARWIN_AMD64/* $AXIS_PATH/build/geropkg/czero/lib/
      fi
      cd $BUILD_PATH

      if [ $os == "windows-amd64" ];then
        rm -rf ./gaxis-*-$os.zip
        zip -r gaxis-$version-$os.zip geropkg/*
      else
         rm -rf ./gaxis-*-$os.tar.gz
         tar czvf gaxis-$version-$os.tar.gz geropkg/*
      fi

      cd $LOCAL_PATH

    done
rm -rf $BUILD_PATH/geropkg/bin
rm -rf $BUILD_PATH/geropkg/czero

