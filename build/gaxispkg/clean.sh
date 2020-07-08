#!/bin/sh

ROOT=$(cd `dirname $0`; pwd)

DATADIR="${ROOT}/data"
if [ ! -z "$1" ]; then
    DATADIR=$1
fi

sh ${ROOT}/stop.sh

echo "rm -rf ${DATADIR}/gaxis/chaindata"
rm -rf ${DATADIR}/gaxis/chaindata
echo "rm -rf ${DATADIR}/gaxis.ipc"
rm -rf ${DATADIR}/gaxis.ipc
echo "rm -rf ${DATADIR}/balance"
rm -rf ${DATADIR}/balance
echo "rm -rf ${DATADIR}/exchange"
rm -rf ${DATADIR}/exchange
echo "rm -rf ${DATADIR}/stake"
rm -rf ${DATADIR}/stake
echo "rm -rf ${DATADIR}/light"
rm -rf ${DATADIR}/light
