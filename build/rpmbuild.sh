#!/bin/bash
START_PATH=`pwd`
SOURCE_PATH=`pwd`/../../
NAME_PATTERN="Name:"
VERSION_PATTERN="Version:"
while read line; do
    if [[ $line =~ ${NAME_PATTERN} ]]; then
        fields=(`echo $line`);
        NAME=${fields[1]}
        echo "NAME:${NAME}"
    elif [[ $line =~ ${VERSION_PATTERN} ]]; then
        fields=(`echo $line`);
        VERSION=${fields[1]}
        echo "VERSION:${VERSION}"
    fi
done <${START_PATH}/gaxis.spec

if [ ! -d ./package ]; then
    mkdir ./package
else
    rm -rf ./package
    mkdir ./package
fi

if [ ! -f ${START_PATH}/bin/gaxis ]; then
    cd ${START_PATH}/../
    make clean 
    make all
fi


if [ ! -d ${SOURCE_PATH}/go-axis-import ]; then
    echo "there is no project available for package"
    exit 1
fi

cd ${START_PATH}/package
mkdir BUILD BUILDROOT RPMS SPECS SRPMS SOURCES
cp ${START_PATH}/gaxis.spec SPECS
cp ${START_PATH}/gaxis.conf 
cd ${START_PATH}/
if [ ! -d ${START_PATH}/tmp ]; then
    mkdir -p ${START_PATH}/tmp/${NAME}-${VERSION}
else
    rm -rf ${START_PATH}/tmp
    mkdir -p ${START_PATH}/tmp/${NAME}-${VERSION}
fi
cd ${START_PATH}/tmp/${NAME}-${VERSION}
mkdir -p etc/gaxis
mkdir -p usr/local/bin
mkdir -p usr/lib64
cp -rf ${START_PATH}/bin/* usr/local/bin
mv  usr/local/bin/gaxis usr/local/bin/_gero
cp -rf ${SOURCE_PATH}/go-axis-import/czero* usr/lib64
cat > usr/local/bin/gaxis << EOL
export LD_LIBRARY_PATH=/usr/lib64/czero/lib
/usr/local/bin/_gero $1
EOL
#echo "export LD_LIBRARY_PATH=/usr/lib64/czero/lib\n/usr/bin/_gero $1\n" > usr/bin/gaxis
cp ${START_PATH}/gaxis.conf etc/gaxis/
chmod -R 755 *
cd ${START_PATH}/tmp
tar -czvf ${NAME}-${VERSION}.tar.gz ./${NAME}-${VERSION}
cd ${START_PATH}/package
cp ${START_PATH}/tmp/${NAME}-${VERSION}.tar.gz SOURCES/

cd ${START_PATH}/package
rpmbuild -ba SPECS/gaxis.spec
